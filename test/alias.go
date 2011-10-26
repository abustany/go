// errchk $G -e $D/$F.go

// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Test that error messages say what the source file says
// (uint8 vs byte).

import (
	"fmt"
	"utf8"
)

func f(byte) {}
func g(uint8) {}

func main() {
	var x float64
	f(x)  // ERROR "byte"
	g(x)  // ERROR "uint8"

	// Test across imports.

	var ff fmt.Formatter
	var fs fmt.State
	ff.Format(fs, x)  // ERROR "rune"

	utf8.RuneStart(x)  // ERROR "byte"

	var s utf8.String
	s.At(x)  // ERROR "int"
}
