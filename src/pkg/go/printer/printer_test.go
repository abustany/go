// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package printer

import (
	"bytes";
	"flag";
	"io";
	"go/ast";
	"go/parser";
	"go/printer";
	"os";
	"path";
	"tabwriter";
	"testing";
)


const (
	dataDir = "testdata";
	tabwidth = 4;
	padding = 1;
	tabchar = '\t';
)


var update = flag.Bool("update", false, "update golden files");


func lineString(text []byte, i int) string {
	i0 := i;
	for i < len(text) && text[i] != '\n' {
		i++;
	}
	return string(text[i0 : i]);
}


func check(t *testing.T, source, golden string, exports bool) {
	// parse source
	prog, err := parser.ParseFile(source, nil, parser.ParseComments);
	if err != nil {
		t.Error(err);
		return;
	}

	// filter exports if necessary
	if exports {
		ast.FilterExports(prog);  // ignore result
	}

	// format source
	var buf bytes.Buffer;
	w := tabwriter.NewWriter(&buf, tabwidth, padding, tabchar, 0);
	Fprint(w, prog, DocComments);
	w.Flush();
	res := buf.Data();

	// update golden files if necessary
	if *update {
		if err := io.WriteFile(golden, res, 0644); err != nil {
			t.Error(err);
		}
		return;
	}

	// get golden
	gld, err := io.ReadFile(golden);
	if err != nil {
		t.Error(err);
		return;
	}

	// compare lengths
	if len(res) != len(gld) {
		t.Errorf("len = %d, expected %d (= len(%s))", len(res), len(gld), golden);
	}

	// compare contents
	for i, line, offs := 0, 1, 0; i < len(res) && i < len(gld); i++ {
		ch := res[i];
		if ch != gld[i] {
			t.Errorf("%s:%d:%d: %s", source, line, i-offs+1, lineString(res, offs));
			t.Errorf("%s:%d:%d: %s", golden, line, i-offs+1, lineString(gld, offs));
			t.Error();
			return;
		}
		if ch == '\n' {
			line++;
			offs = i+1;
		}
	}
}


type entry struct {
	source, golden string;
	exports bool;
}

// Use gotest -update to create/update the respective golden files.
var data = []entry{
	entry{ "source1.go", "golden1.go", false },
	entry{ "source1.go", "golden1.x", true },
}


func Test(t *testing.T) {
	for _, e := range data {
		check(t, path.Join(dataDir, e.source), path.Join(dataDir, e.golden), e.exports);
	}
}
