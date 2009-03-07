// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package contains an interface to the low-level operating system
// primitives.  The details vary depending on the underlying system.
// Its primary use is inside other packages that provide a more portable
// interface to the system, such as "os", "time" and "net".  Use those
// packages rather than this one if you can.
// For details of the functions and data types in this package consult
// the manuals for the appropriate operating system.
package syscall

/*
 * Foundation of system call interface.
 */

func Syscall(trap int64, a1, a2, a3 int64) (r1, r2, err int64);
func Syscall6(trap int64, a1, a2, a3, a4, a5, a6 int64) (r1, r2, err int64);
func RawSyscall(trap int64, a1, a2, a3 int64) (r1, r2, err int64);

/*
 * Used to convert file names to byte arrays for passing to kernel,
 * but useful elsewhere too.
 */
func StringBytePtr(s string) *byte {
	a := make([]byte, len(s)+1);
	for i := 0; i < len(s); i++ {
		a[i] = s[i];
	}
	return &a[0];
}
