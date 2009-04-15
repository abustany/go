// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Template library.  See http://code.google.com/p/json-template/wiki/Reference
// TODO: document this here as well.
package template

import (
	"fmt";
	"io";
	"os";
	"reflect";
	"strings";
	"template";
)

var ErrUnmatchedRDelim = os.NewError("unmatched closing delimiter")
var ErrUnmatchedLDelim = os.NewError("unmatched opening delimiter")
var ErrBadDirective = os.NewError("unrecognized directive name")
var ErrEmptyDirective = os.NewError("empty directive")
var ErrFields = os.NewError("incorrect fields for directive")
var ErrSyntax = os.NewError("directive out of place")
var ErrNoEnd = os.NewError("section does not have .end")
var ErrNoVar = os.NewError("variable name not in struct");
var ErrBadType = os.NewError("unsupported type for variable");
var ErrNotStruct = os.NewError("driver must be a struct")
var ErrNoFormatter = os.NewError("unknown formatter")
var ErrEmptyDelims = os.NewError("empty delimiter strings")

// All the literals are aces.
var lbrace = []byte{ '{' }
var rbrace = []byte{ '}' }
var space = []byte{ ' ' }

// The various types of "tokens", which are plain text or (usually) brace-delimited descriptors
const (
	Alternates = iota;
	Comment;
	End;
	Literal;
	Or;
	Repeated;
	Section;
	Text;
	Variable;
)

// FormatterMap is the type describing the mapping from formatter
// names to the functions that implement them.
type FormatterMap map[string] func(io.Write, interface{}, string)

// Built-in formatters.
var builtins = FormatterMap {
	"html" : HtmlFormatter,
	"str" : StringFormatter,
	"" : StringFormatter,
}

// State for executing a Template
type state struct {
	parent	*state;	// parent in hierarchy
	errorchan	chan *os.Error;	// for erroring out
	data	reflect.Value;	// the driver data for this section etc.
	wr	io.Write;	// where to send output
}

// Report error and stop generation.
func (st *state) error(err *os.Error, args ...) {
	fmt.Fprintf(os.Stderr, "template: %v%s\n", err, fmt.Sprint(args));
	st.errorchan <- err;
	sys.Goexit();
}

type Template struct {
	fmap	FormatterMap;	// formatters for variables
	ldelim, rdelim	[]byte;	// delimiters; default {}
	buf	[]byte;	// input text to process
	p	int;	// position in buf
	linenum	*int;	// position in input
}

// Initialize a top-level template in prepratation for parsing.
// The formatter map and delimiters are already set.
func (t *Template) init(buf []byte) *Template {
	t.buf = buf;
	t.p = 0;
	t.linenum = new(int);
	return t;
}
// Create a template deriving from its parent
func childTemplate(parent *Template, buf []byte) *Template {
	t := new(Template);
	t.ldelim = parent.ldelim;
	t.rdelim = parent.rdelim;
	t.buf = buf;
	t.p = 0;
	t.fmap = parent.fmap;
	t.linenum = parent.linenum;
	return t;
}

func white(c uint8) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// safely, does s[n:n+len(t)] == t?
func equal(s []byte, n int, t []byte) bool {
	b := s[n:len(s)];
	if len(t) > len(b) {	// not enough space left for a match.
		return false
	}
	for i , c := range t {
		if c != b[i] {
			return false
		}
	}
	return true
}

func (t *Template) execute(st *state)
func (t *Template) executeSection(w []string, st *state)

// nextItem returns the next item from the input buffer.  If the returned
// item is empty, we are at EOF.  The item will be either a
// delimited string or a non-empty string between delimited
// strings.  Most tokens stop at (but include, if plain text) a newline.
// Action tokens on a line by themselves drop the white space on
// either side, up to and including the newline.
func (t *Template) nextItem(st *state) []byte {
	sawLeft := false;	// are we waiting for an opening delimiter?
	special := false;	// is this a {.foo} directive, which means trim white space?
	// Delete surrounding white space if this {.foo} is the only thing on the line.
	trim_white := t.p == 0 || t.buf[t.p-1] == '\n';
	only_white := true;	// we have seen only white space so far
	var i int;
	start := t.p;
Loop:
	for i = t.p; i < len(t.buf); i++ {
		switch {
		case t.buf[i] == '\n':
			*t.linenum++;
			i++;
			break Loop;
		case white(t.buf[i]):
			// white space, do nothing
		case !sawLeft && equal(t.buf, i, t.ldelim):  // sawLeft checked because delims may be equal
			// anything interesting already on the line?
			if !only_white {
				break Loop;
			}
			// is it a directive or comment?
			j := i + len(t.ldelim);  // position after delimiter
			if j+1 < len(t.buf) && (t.buf[j] == '.' || t.buf[j] == '#') {
				special = true;
				if trim_white && only_white {
					start = i;
				}
			} else if i > t.p {  // have some text accumulated so stop before delimiter
				break Loop;
			}
			sawLeft = true;
			i = j - 1;
		case equal(t.buf, i, t.rdelim):
			if !sawLeft {
				st.error(ErrUnmatchedRDelim)
			}
			sawLeft = false;
			i += len(t.rdelim);
			break Loop;
		default:
			only_white = false;
		}
	}
	if sawLeft {
		st.error(ErrUnmatchedLDelim)
	}
	item := t.buf[start:i];
	if special && trim_white {
		// consume trailing white space
		for ; i < len(t.buf) && white(t.buf[i]); i++ {
			if t.buf[i] == '\n' {
				i++;
				break	// stop after newline
			}
		}
	}
	t.p = i;
	return item
}

// Turn a byte array into a white-space-split array of strings.
func words(buf []byte) []string {
	s := make([]string, 0, 5);
	p := 0; // position in buf
	// one word per loop
	for i := 0; ; i++ {
		// skip white space
		for ; p < len(buf) && white(buf[p]); p++ {
		}
		// grab word
		start := p;
		for ; p < len(buf) && !white(buf[p]); p++ {
		}
		if start == p {	// no text left
			break
		}
		if i == cap(s) {
			ns := make([]string, 2*cap(s));
			for j := range s {
				ns[j] = s[j]
			}
			s = ns;
		}
		s = s[0:i+1];
		s[i] = string(buf[start:p])
	}
	return s
}

// Analyze an item and return its type and, if it's an action item, an array of
// its constituent words.
func (t *Template) analyze(item []byte, st *state) (tok int, w []string) {
	// item is known to be non-empty
	if !equal(item, 0, t.ldelim) {	// doesn't start with left delimiter
		tok = Text;
		return
	}
	if !equal(item, len(item)-len(t.rdelim), t.rdelim) {	// doesn't end with right delimiter
		st.error(ErrUnmatchedLDelim)  // should not happen anyway
	}
	if len(item) <= len(t.ldelim)+len(t.rdelim) {	// no contents
		st.error(ErrEmptyDirective)
	}
	// Comment
	if item[len(t.ldelim)] == '#' {
		tok = Comment;
		return
	}
	// Split into words
	w = words(item[len(t.ldelim): len(item)-len(t.rdelim)]);  // drop final delimiter
	if len(w) == 0 {
		st.error(ErrBadDirective)
	}
	if len(w[0]) == 0 {
		st.error(ErrEmptyDirective)
	}
	if len(w) == 1 && w[0][0] != '.' {
		tok = Variable;
		return;
	}
	switch w[0] {
	case ".meta-left", ".meta-right", ".space":
		tok = Literal;
		return;
	case ".or":
		tok = Or;
		return;
	case ".end":
		tok = End;
		return;
	case ".section":
		if len(w) != 2 {
			st.error(ErrFields, ": ", string(item))
		}
		tok = Section;
		return;
	case ".repeated":
		if len(w) != 3 || w[1] != "section" {
			st.error(ErrFields, ": ", string(item))
		}
		tok = Repeated;
		return;
	case ".alternates":
		if len(w) != 2 || w[1] != "with" {
			st.error(ErrFields, ": ", string(item))
		}
		tok = Alternates;
		return;
	}
	st.error(ErrBadDirective, ": ", string(item));
	return
}

// If the data for this template is a struct, find the named variable.
// The special name "@" denotes the current data.
func (st *state) findVar(s string) reflect.Value {
	if s == "@" {
		return st.data
	}
	data := reflect.Indirect(st.data);
	typ, ok := data.Type().(reflect.StructType);
	if ok {
		for i := 0; i < typ.Len(); i++ {
			name, ftyp, tag, offset := typ.Field(i);
			if name == s {
				return data.(reflect.StructValue).Field(i)
			}
		}
	}
	return nil
}

// Is there no data to look at?
func empty(v reflect.Value, indirect_ok bool) bool {
	v = reflect.Indirect(v);
	if v == nil {
		return true
	}
	switch v.Type().Kind() {
	case reflect.StructKind:
		return false;
	case reflect.ArrayKind:
		return v.(reflect.ArrayValue).Len() == 0;
	}
	return true;
}

// Execute a ".repeated" section
func (t *Template) executeRepeated(w []string, st *state) {
	if w[1] != "section" {
		st.error(ErrSyntax, `: .repeated must have "section"`)
	}

	// Find driver array/struct for this section.  It must be in the current struct.
	field := st.findVar(w[2]);
	if field == nil {
		st.error(ErrNoVar, ": .repeated ", w[2], " in ", reflect.Indirect(st.data).Type());
	}

	// Must be an array/slice
	if field != nil && field.Kind() != reflect.ArrayKind {
		st.error(ErrBadType, " in .repeated: ", w[2], " ", field.Type().String());
	}
	// Scan repeated section, remembering slice of text we must execute.
	nesting := 0;
	start := t.p;
	end := t.p;
Loop:
	for {
		item := t.nextItem(st);
		if len(item) ==  0 {
			st.error(ErrNoEnd)
		}
		tok, s := t.analyze(item, st);
		switch tok {
		case Comment:
			continue;	// just ignore it
		case End:
			if nesting == 0 {
				break Loop
			}
			nesting--;
		case Repeated, Section:
			nesting++;
		case Literal, Or, Text, Variable:
			// just accumulate
		default:
			panic("unknown section item", string(item));
		}
		end = t.p
	}
	if field != nil {
		array := field.(reflect.ArrayValue);
		for i := 0; i < array.Len(); i++ {
			tmp := childTemplate(t, t.buf[start:end]);
			tmp.execute(&state{st, st.errorchan, array.Elem(i), st.wr});
		}
	}
}

// Execute a ".section"
func (t *Template) executeSection(w []string, st *state) {
	// Find driver data for this section.  It must be in the current struct.
	field := st.findVar(w[1]);
	if field == nil {
		st.error(ErrNoVar, ": .section ", w[1], " in ", reflect.Indirect(st.data).Type());
	}
	// Scan section, remembering slice of text we must execute.
	orFound := false;
	nesting := 0;  // How deeply are .section and .repeated nested?
	start := t.p;
	end := t.p;
	accumulate := !empty(field, true);	// Keep this section if there's data
Loop:
	for {
		item := t.nextItem(st);
		if len(item) ==  0 {
			st.error(ErrNoEnd)
		}
		tok, s := t.analyze(item, st);
		switch tok {
		case Comment:
			continue;	// just ignore it
		case End:
			if nesting == 0 {
				break Loop
			}
			nesting--;
		case Or:
			if nesting > 0 {	// just accumulate
				break
			}
			if orFound {
				st.error(ErrSyntax, ": .or");
			}
			orFound = true;
			if !accumulate {
				// No data; execute the .or instead
				start = t.p;
				end = t.p;
				accumulate = true;
				continue;
			} else {
				// Data present so disregard the .or section
				accumulate = false
			}
		case Repeated, Section:
			nesting++;
		case Literal, Text, Variable:
			// just accumulate
		default:
			panic("unknown section item", string(item));
		}
		if accumulate {
			end = t.p
		}
	}
	tmp := childTemplate(t, t.buf[start:end]);
	tmp.execute(&state{st, st.errorchan, field, st.wr});
}

// Look up a variable, up through the parent if necessary.
func (t *Template) varValue(name string, st *state) reflect.Value {
	field := st.findVar(name);
	if field == nil {
		if st.parent == nil {
			st.error(ErrNoVar, ": ", name)
		}
		return t.varValue(name, st.parent);
	}
	return field;
}

// Evaluate a variable, looking up through the parent if necessary.
// If it has a formatter attached ({var|formatter}) run that too.
func (t *Template) writeVariable(st *state, name_formatter string) {
	name := name_formatter;
	formatter := "";
	bar := strings.Index(name_formatter, "|");
	if bar >= 0 {
		name = name_formatter[0:bar];
		formatter = name_formatter[bar+1:len(name_formatter)];
	}
	val := t.varValue(name, st).Interface();
	// is it in user-supplied map?
	if t.fmap != nil {
		if fn, ok := t.fmap[formatter]; ok {
			fn(st.wr, val, formatter);
			return;
		}
	}
	// is it in builtin map?
	if fn, ok := builtins[formatter]; ok {
		fn(st.wr, val, formatter);
		return;
	}
	st.error(ErrNoFormatter, ": ", formatter);
	panic("notreached");
}

func (t *Template) execute(st *state) {
	for {
		item := t.nextItem(st);
		if len(item) == 0 {
			return
		}
		tok, w := t.analyze(item, st);
		switch tok {
		case Comment:
			break;
		case Text:
			st.wr.Write(item);
		case Literal:
			switch w[0] {
			case ".meta-left":
				st.wr.Write(t.ldelim);
			case ".meta-right":
				st.wr.Write(t.rdelim);
			case ".space":
				st.wr.Write(space);
			default:
				panic("unknown literal: ", w[0]);
			}
		case Variable:
			t.writeVariable(st, w[0]);
		case Or, End, Alternates:
			st.error(ErrSyntax, ": ", string(item));
		case Section:
			t.executeSection(w, st);
		case Repeated:
			t.executeRepeated(w, st);
		default:
			panic("bad directive in execute:", string(item));
		}
	}
}

func (t *Template) doParse() {
	// stub for now
}

// Parse initializes a Template by parsing its definition.  The string s contains
// the template text.  If any errors occur, it returns the error and line number
// in the text of the erroneous construct.
func (t *Template) Parse(s string) (*os.Error, int) {
	if len(t.ldelim) == 0 || len(t.rdelim) == 0 {
		return ErrEmptyDelims, 0
	}
	t.init(io.StringBytes(s));
	ch := make(chan *os.Error);
	go func() {
		t.doParse();
		ch <- nil;	// clean return;
	}();
	err := <-ch;
	if err != nil {
		return err, *t.linenum
	}
	return nil, 0
}

// Execute executes a parsed template on the specified data object,
// generating output to wr.
func (t *Template) Execute(data interface{}, wr io.Write) *os.Error {
	// Extract the driver data.
	val := reflect.NewValue(data);
	ch := make(chan *os.Error);
	go func() {
		t.p = 0;
		t.execute(&state{nil, ch, val, wr});
		ch <- nil;	// clean return;
	}();
	return <-ch;
}

// New creates a new template with the specified formatter map (which
// may be nil) defining auxiliary functions for formatting variables.
func New(fmap FormatterMap) *Template {
	t := new(Template);
	t.fmap = fmap;
	t.ldelim = lbrace;
	t.rdelim = rbrace;
	return t;
}

// SetDelims sets the left and right delimiters for operations in the template.
func (t *Template) SetDelims(left, right string) {
	t.ldelim = io.StringBytes(left);
	t.rdelim = io.StringBytes(right);
}

// Parse creates a Template with default parameters (such as {} for
// metacharacters).  The string s contains the template text and the
// formatter map fmap (which may be nil) defines auxiliary functions
// for formatting variables.  It returns the template, an error report
// (or nil), and the line number in the text of the erroneous construct.
func Parse(s string, fmap FormatterMap) (*Template, *os.Error, int) {
	t := New(fmap);
	err, line := t.Parse(s);
	return t, err, line
}
