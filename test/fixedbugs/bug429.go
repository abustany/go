// $G $D/$F.go && $L $F.$A && ! ./$A.out || echo BUG: bug429

// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Should print deadlock message, not hang.

package main

func main() {
	select{}
}
