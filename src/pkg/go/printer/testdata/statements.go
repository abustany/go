// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package statements

var expr bool;

func use(x interface{}) {}

// Formatting of if-statement headers.
func _() {
	if {}
	if;{}  // no semicolon printed
	if expr{}
	if;expr{}  // no semicolon printed
	if x:=expr;{
	use(x)}
	if x:=expr; expr {use(x)}
}


// Formatting of switch-statement headers.
func _() {
	switch {}
	switch;{}  // no semicolon printed
	switch expr {}
	switch;expr{}  // no semicolon printed
	switch x := expr; { default:use(
x)
	}
	switch x := expr; expr {default:use(x)}
}


// Formatting of switch statement bodies.
func _() {
	switch {
	}

	switch x := 0; x {
	case 1:
		use(x);
		use(x);  // followed by an empty line

	case 2:  // followed by an empty line

		use(x);  // followed by an empty line

	case 3:  // no empty lines
		use(x);
		use(x);
	}
}


// Formatting of for-statement headers.
func _() {
	for{}
	for expr {}
	for;;{}  // no semicolon printed
	for x :=expr;; {use( x)}
	for; expr;{}  // no semicolon printed
	for; ; expr = false {}
	for x :=expr; expr; {use(x)}
	for x := expr;; expr=false {use(x)}
	for;expr;expr =false {
	}
	for x := expr;expr;expr = false { use(x) }
	for x := range []int{} { use(x) }
}


// Extra empty lines inside functions. Do respect source code line
// breaks between statement boundaries but print at most one empty
// line at a time.
func _() {

	const _ = 0;

	const _ = 1;
	type _ int;
	type _ float;

	var _ = 0;
	var x = 1;

	// Each use(x) call below should have at most one empty line before and after.



	use(x);

	if x < x {

		use(x);

	} else {

		use(x);

	}
}
