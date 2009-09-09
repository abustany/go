// errchk $G $D/$F.go

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var t []int
var s string;
var m map[string]int;

func main() {
	println(t["hi"]);	// ERROR "non-integer"
	println(s["hi"]);	// ERROR "non-integer"
	println(m[0]);	// ERROR "map index"
}

