// $G $D/$F.go || echo BUG: bug173

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// these used to fail because the runtime
// functions that get called to implement them
// expected string, not T.

package main

type T string
func main() {
	var t T = "hello";
	println(t[0:4], t[4]);
	for i, x := range t {
	}
	for i := range t {
	}
}
