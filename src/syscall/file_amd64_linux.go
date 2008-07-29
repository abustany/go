// Copyright 2009 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package syscall

import syscall "syscall"

export Stat
export stat, fstat, lstat
export open, close, read, write, pipe

func	StatToInt(s *Stat) int64;
func	Addr32ToInt(s *int32) int64;

// Stat and relatives for Linux

type dev_t uint64;
type ino_t uint64;
type mode_t uint32;
type nlink_t uint64;
type uid_t uint32;
type gid_t uint32;
type off_t int64;
type blksize_t int64;
type blkcnt_t int64;
type time_t int64;

type Timespec struct {
	tv_sec	time_t;
	tv_nsec	int64;
}

type Stat struct {
	st_dev	dev_t;     /* ID of device containing file */
	st_ino	ino_t;     /* inode number */
	st_nlink	nlink_t;   /* number of hard links */
	st_mode	mode_t;    /* protection */
	st_uid	uid_t;     /* user ID of owner */
	st_gid	gid_t;     /* group ID of owner */
	pad0	int32;
	st_rdev	dev_t;    /* device ID (if special file) */
	st_size	off_t;    /* total size, in bytes */
	st_blksize	blksize_t; /* blocksize for filesystem I/O */
	st_blocks	blkcnt_t;  /* number of blocks allocated */
	st_atime	Timespec;   /* time of last access */
	st_mtime	Timespec;   /* time of last modification */
	st_ctime	Timespec;   /* time of last status change */
	st_unused4	int64;
	st_unused5	int64;
	st_unused6	int64;
}

func open(name *byte, mode int64) (ret int64, errno int64) {
	const SYSOPEN = 2;
	r1, r2, err := syscall.Syscall(SYSOPEN, AddrToInt(name), mode, 0);
	return r1, err;
}

func close(fd int64) (ret int64, errno int64) {
	const SYSCLOSE = 3;
	r1, r2, err := syscall.Syscall(SYSCLOSE, fd, 0, 0);
	return r1, err;
}

func read(fd int64, buf *byte, nbytes int64) (ret int64, errno int64) {
print "READ: ", fd, " ", nbytes, "\n";
	const SYSREAD = 0;
	r1, r2, err := syscall.Syscall(SYSREAD, fd, AddrToInt(buf), nbytes);
	return r1, err;
}

func write(fd int64, buf *byte, nbytes int64) (ret int64, errno int64) {
	const SYSWRITE = 1;
	r1, r2, err := syscall.Syscall(SYSWRITE, fd, AddrToInt(buf), nbytes);
	return r1, err;
}

func pipe(fds *[2]int64) (ret int64, errno int64) {
	const SYSPIPE = 22;
	var t [2] int32;
	r1, r2, err := syscall.Syscall(SYSPIPE, Addr32ToInt(&t[0]), 0, 0);
	if r1 < 0 {
		return r1, err;
	}
	fds[0] = int64(t[0]);
	fds[1] = int64(t[1]);
	return 0, err;
}

func stat(name *byte, buf *Stat) (ret int64, errno int64) {
	const SYSSTAT = 4;
	r1, r2, err := syscall.Syscall(SYSSTAT, AddrToInt(name), StatToInt(buf), 0);
	return r1, err;
}

func lstat(name *byte, buf *Stat) (ret int64, errno int64) {
	const SYSLSTAT = 6;
	r1, r2, err := syscall.Syscall(SYSLSTAT, AddrToInt(name), StatToInt(buf), 0);
	return r1, err;
}

func fstat(fd int64, buf *Stat) (ret int64, errno int64) {
	const SYSFSTAT = 5;
	r1, r2, err := syscall.Syscall(SYSFSTAT, fd, StatToInt(buf), 0);
	return r1, err;
}

