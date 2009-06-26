// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import (
	"os";
	"syscall";
)

// ForkExec forks the current process and invokes Exec with the file, arguments,
// and environment specified by argv0, argv, and envv.  It returns the process
// id of the forked process and an Error, if any.  The fd array specifies the
// file descriptors to be set up in the new process: fd[0] will be Unix file
// descriptor 0 (standard input), fd[1] descriptor 1, and so on.  A nil entry
// will cause the child to have no open file descriptor with that index.
// If dir is not empty, the child chdirs into the directory before execing the program.
func ForkExec(argv0 string, argv []string, envv []string, dir string, fd []*File)
	(pid int, err Error)
{
	// Create array of integer (system) fds.
	intfd := make([]int, len(fd));
	for i, f := range fd {
		if f == nil {
			intfd[i] = -1;
		} else {
			intfd[i] = f.Fd();
		}
	}

	p, e := syscall.ForkExec(argv0, argv, envv, dir, intfd);
	if e != 0 {
		return 0, &PathError{"fork/exec", argv0, Errno(e)};
	}
	return p, nil;
}

// Exec replaces the current process with an execution of the program
// named by argv0, with arguments argv and environment envv.
// If successful, Exec never returns.  If it fails, it returns an Error.
// ForkExec is almost always a better way to execute a program.
func Exec(argv0 string, argv []string, envv []string) Error {
	if envv == nil {
		envv = Environ();
	}
	e := syscall.Exec(argv0, argv, envv);
	if e != 0 {
		return &PathError{"exec", argv0, Errno(e)};
	}
	return nil;
}

// TODO(rsc): Should os implement its own syscall.WaitStatus
// wrapper with the methods, or is exposing the underlying one enough?
//
// TODO(rsc): Certainly need to have os.Rusage struct,
// since syscall one might have different field types across
// different OS.

// Waitmsg stores the information about an exited process as reported by Wait.
type Waitmsg struct {
	Pid int;	// The process's id.
	syscall.WaitStatus;	// System-dependent status info.
	Rusage *syscall.Rusage;	// System-dependent resource usage info.
}

// Options for Wait.
const (
	WNOHANG = syscall.WNOHANG;	// Don't wait if no process has exited.
	WSTOPPED = syscall.WSTOPPED;	// If set, status of stopped subprocesses is also reported.
	WUNTRACED = WSTOPPED;
	WRUSAGE = 1<<30;	// Record resource usage.
)

// Wait waits for process pid to exit or stop, and then returns a
// Waitmsg describing its status and an Error, if any. The options
// (WNOHANG etc.) affect the behavior of the Wait call.
func Wait(pid int, options int) (w *Waitmsg, err Error) {
	var status syscall.WaitStatus;
	var rusage *syscall.Rusage;
	if options & WRUSAGE != 0 {
		rusage = new(syscall.Rusage);
		options ^= WRUSAGE;
	}
	pid1, e := syscall.Wait4(pid, &status, options, rusage);
	if e != 0 {
		return nil, NewSyscallError("wait", e);
	}
	w = new(Waitmsg);
	w.Pid = pid1;
	w.WaitStatus = status;
	w.Rusage = rusage;
	return w, nil;
}

// Getpid returns the process id of the caller.
func Getpid() int {
	p, r2, e := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0);
	return int(p)
}

// Getppid returns the process id of the caller's parent.
func Getppid() int {
	p, r2, e := syscall.Syscall(syscall.SYS_GETPPID, 0, 0, 0);
	return int(p)
}
