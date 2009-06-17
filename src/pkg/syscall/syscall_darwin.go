// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Darwin system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and
// wrap it in our own nicer implementation.

package syscall

import (
	"syscall";
	"unsafe";
)

const OS = "darwin"

/*
 * Pseudo-system calls
 */
// The const provides a compile-time constant so clients
// can adjust to whether there is a working Getwd and avoid
// even linking this function into the binary.  See ../os/getwd.go.
const ImplementsGetwd = false

func Getwd() (string, int) {
	return "", ENOTSUP;
}


/*
 * Wrapped
 */

//sys	getgroups(ngid int, gid *_Gid_t) (n int, errno int)
//sys	setgroups(ngid int, gid *_Gid_t) (errno int)

func Getgroups() (gids []int, errno int) {
	n, err := getgroups(0, nil);
	if err != 0 {
		return nil, errno;
	}
	if n == 0 {
		return nil, 0;
	}

	// Sanity check group count.  Max is 16 on BSD.
	if n < 0 || n > 1000 {
		return nil, EINVAL;
	}

	a := make([]_Gid_t, n);
	n, err = getgroups(n, &a[0]);
	if err != 0 {
		return nil, errno;
	}
	gids = make([]int, n);
	for i, v := range a[0:n] {
		gids[i] = int(v);
	}
	return;
}

func Setgroups(gids []int) (errno int) {
	if len(gids) == 0 {
		return setgroups(0, nil);
	}

	a := make([]_Gid_t, len(gids));
	for i, v := range gids {
		a[i] = _Gid_t(v);
	}
	return setgroups(len(a), &a[0]);
}

// Wait status is 7 bits at bottom, either 0 (exited),
// 0x7F (stopped), or a signal number that caused an exit.
// The 0x80 bit is whether there was a core dump.
// An extra number (exit code, signal causing a stop)
// is in the high bits.

type WaitStatus uint32

const (
	mask = 0x7F;
	core = 0x80;
	shift = 8;

	exited = 0;
	stopped = 0x7F;
)

func (w WaitStatus) Exited() bool {
	return w&mask == exited;
}

func (w WaitStatus) ExitStatus() int {
	if w&mask != exited {
		return -1;
	}
	return int(w >> shift);
}

func (w WaitStatus) Signaled() bool {
	return w&mask != stopped && w&mask != 0;
}

func (w WaitStatus) Signal() int {
	sig := int(w & mask);
	if sig == stopped || sig == 0 {
		return -1;
	}
	return sig;
}

func (w WaitStatus) CoreDump() bool {
	return w.Signaled() && w&core != 0;
}

func (w WaitStatus) Stopped() bool {
	return w&mask == stopped && w>>shift != SIGSTOP;
}

func (w WaitStatus) Continued() bool {
	return w&mask == stopped && w>>shift == SIGSTOP;
}

func (w WaitStatus) StopSignal() int {
	if !w.Stopped() {
		return -1;
	}
	return int(w >> shift) & 0xFF;
}

//sys	wait4(pid int, wstatus *_C_int, options int, rusage *Rusage) (wpid int, errno int)
func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, errno int) {
	var status _C_int;
	wpid, errno = wait4(pid, &status, options, rusage);
	if wstatus != nil {
		*wstatus = WaitStatus(status);
	}
	return;
}

//sys	pipe() (r int, w int, errno int)
func Pipe(p []int) (errno int) {
	if len(p) != 2 {
		return EINVAL;
	}
	p[0], p[1], errno = pipe();
	return;
}

func Sleep(ns int64) (errno int) {
	tv := NsecToTimeval(ns);
	return Select(0, nil, nil, nil, &tv);
}

//sys	accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, errno int)
//sys	bind(s int, addr uintptr, addrlen _Socklen) (errno int)
//sys	connect(s int, addr uintptr, addrlen _Socklen) (errno int)
//sys	socket(domain int, typ int, proto int) (fd int, errno int)
//sys	setsockopt(s int, level int, name int, val uintptr, vallen int) (errno int)

// For testing: clients can set this flag to force
// creation of IPv6 sockets to return EAFNOSUPPORT.
var SocketDisableIPv6 bool

type Sockaddr interface {
	sockaddr() (ptr uintptr, len _Socklen, errno int);	// lowercase; only we can define Sockaddrs
}

type SockaddrInet4 struct {
	Port int;
	Addr [4]byte;
	raw RawSockaddrInet4;
}

func (sa *SockaddrInet4) sockaddr() (uintptr, _Socklen, int) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return 0, 0, EINVAL;
	}
	sa.raw.Len = SizeofSockaddrInet4;
	sa.raw.Family = AF_INET;
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port));
	p[0] = byte(sa.Port>>8);
	p[1] = byte(sa.Port);
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i];
	}
	return uintptr(unsafe.Pointer(&sa.raw)), _Socklen(sa.raw.Len), 0;
}

type SockaddrInet6 struct {
	Port int;
	Addr [16]byte;
	raw RawSockaddrInet6;
}

func (sa *SockaddrInet6) sockaddr() (uintptr, _Socklen, int) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return 0, 0, EINVAL;
	}
	sa.raw.Len = SizeofSockaddrInet6;
	sa.raw.Family = AF_INET6;
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port));
	p[0] = byte(sa.Port>>8);
	p[1] = byte(sa.Port);
	for i := 0; i < len(sa.Addr); i++ {
		sa.raw.Addr[i] = sa.Addr[i];
	}
	return uintptr(unsafe.Pointer(&sa.raw)), _Socklen(sa.raw.Len), 0;
}

type SockaddrUnix struct {
	Name string;
	raw RawSockaddrUnix;
}

func (sa *SockaddrUnix) sockaddr() (uintptr, _Socklen, int) {
	name := sa.Name;
	n := len(name);
	if n >= len(sa.raw.Path) || n == 0 {
		return 0, 0, EINVAL;
	}
	sa.raw.Len = byte(3 + n);	// 2 for Family, Len; 1 for NUL
	sa.raw.Family = AF_UNIX;
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i]);
	}
	return uintptr(unsafe.Pointer(&sa.raw)), _Socklen(sa.raw.Len), 0;
}

func anyToSockaddr(rsa *RawSockaddrAny) (Sockaddr, int) {
	switch rsa.Addr.Family {
	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa));
		if pp.Len < 3 || pp.Len > SizeofSockaddrUnix {
			return nil, EINVAL
		}
		sa := new(SockaddrUnix);
		n := int(pp.Len) - 3;	// subtract leading Family, Len, terminating NUL
		for i := 0; i < n; i++ {
			if pp.Path[i] == 0 {
				// found early NUL; assume Len is overestimating
				n = i;
				break;
			}
		}
		bytes := (*[len(pp.Path)]byte)(unsafe.Pointer(&pp.Path[0]));
		sa.Name = string(bytes[0:n]);
		return sa, 0;

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa));
		sa := new(SockaddrInet4);
		p := (*[2]byte)(unsafe.Pointer(&pp.Port));
		sa.Port = int(p[0])<<8 + int(p[1]);
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i];
		}
		return sa, 0;

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa));
		sa := new(SockaddrInet6);
		p := (*[2]byte)(unsafe.Pointer(&pp.Port));
		sa.Port = int(p[0])<<8 + int(p[1]);
		for i := 0; i < len(sa.Addr); i++ {
			sa.Addr[i] = pp.Addr[i];
		}
		return sa, 0;
	}
	return nil, EAFNOSUPPORT;
}

func Accept(fd int) (nfd int, sa Sockaddr, errno int) {
	var rsa RawSockaddrAny;
	var len _Socklen = SizeofSockaddrAny;
	nfd, errno = accept(fd, &rsa, &len);
	if errno != 0 {
		return;
	}
	sa, errno = anyToSockaddr(&rsa);
	if errno != 0 {
		Close(nfd);
		nfd = 0;
	}
	return;
}

func Bind(fd int, sa Sockaddr) (errno int) {
	ptr, n, err := sa.sockaddr();
	if err != 0 {
		return err;
	}
	return bind(fd, ptr, n);
}

func Connect(fd int, sa Sockaddr) (errno int) {
	ptr, n, err := sa.sockaddr();
	if err != 0 {
		return err;
	}
	return connect(fd, ptr, n);
}

func Socket(domain, typ, proto int) (fd, errno int) {
	if domain == AF_INET6 && SocketDisableIPv6 {
		return -1, EAFNOSUPPORT
	}
	fd, errno = socket(domain, typ, proto);
	return;
}

func SetsockoptInt(fd, level, opt int, value int) (errno int) {
	var n = int32(value);
	return setsockopt(fd, level, opt, uintptr(unsafe.Pointer(&n)), 4);
}

func SetsockoptTimeval(fd, level, opt int, tv *Timeval) (errno int) {
	return setsockopt(fd, level, opt, uintptr(unsafe.Pointer(tv)), unsafe.Sizeof(*tv));
}

func SetsockoptLinger(fd, level, opt int, l *Linger) (errno int) {
	return setsockopt(fd, level, opt, uintptr(unsafe.Pointer(l)), unsafe.Sizeof(*l));
}

//sys	kevent(kq int, change uintptr, nchange int, event uintptr, nevent int, timeout *Timespec) (n int, errno int)
func Kevent(kq int, changes, events []Kevent_t, timeout *Timespec) (n int, errno int) {
	var change, event uintptr;
	if len(changes) > 0 {
		change = uintptr(unsafe.Pointer(&changes[0]));
	}
	if len(events) > 0 {
		event = uintptr(unsafe.Pointer(&events[0]));
	}
	return kevent(kq, change, len(changes), event, len(events), timeout);
}

// TODO: wrap
//	Acct(name nil-string) (errno int)
//	Futimes(fd int, timeval *Timeval) (errno int)	// Pointer to 2 timevals!
//	Gethostuuid(uuid *byte, timeout *Timespec) (errno int)
//	Getpeername(fd int, addr *Sockaddr, addrlen *int) (errno int)
//	Getsockname(fd int, addr *Sockaddr, addrlen *int) (errno int)
//	Getsockopt(s int, level int, name int, val *byte, vallen *int) (errno int)
//	Madvise(addr *byte, len int, behav int) (errno int)
//	Mprotect(addr *byte, len int, prot int) (errno int)
//	Msync(addr *byte, len int, flags int) (errno int)
//	Munmap(addr *byte, len int) (errno int)
//	Ptrace(req int, pid int, addr uintptr, data int) (ret uintptr, errno int)
//	Recvfrom(s int, buf *byte, nbuf int, flags int, from *Sockaddr, fromlen *int) (n int, errno int)
//	Recvmsg(s int, msg *Msghdr, flags int) (n int, errno int)
//	Sendmsg(s int, msg *Msghdr, flags int) (n int, errno int)
//	Sendto(s int, buf *byte, nbuf int, flags int, to *Sockaddr, addrlen int) (errno int)
//	Utimes(path string, timeval *Timeval) (errno int)	// Pointer to 2 timevals!
//sys	fcntl(fd int, cmd int, arg int) (val int, errno int)


/*
 * Exposed directly
 */
//sys	Access(path string, flags int) (errno int)
//sys	Adjtime(delta *Timeval, olddelta *Timeval) (errno int)
//sys	Chdir(path string) (errno int)
//sys	Chflags(path string, flags int) (errno int)
//sys	Chmod(path string, mode int) (errno int)
//sys	Chown(path string, uid int, gid int) (errno int)
//sys	Chroot(path string) (errno int)
//sys	Close(fd int) (errno int)
//sys	Dup(fd int) (nfd int, errno int)
//sys	Dup2(from int, to int) (errno int)
//sys	Exchangedata(path1 string, path2 string, options int) (errno int)
//sys	Exit(code int)
//sys	Fchdir(fd int) (errno int)
//sys	Fchflags(path string, flags int) (errno int)
//sys	Fchmod(fd int, mode int) (errno int)
//sys	Fchown(fd int, uid int, gid int) (errno int)
//sys	Flock(fd int, how int) (errno int)
//sys	Fpathconf(fd int, name int) (val int, errno int)
//sys	Fstat(fd int, stat *Stat_t) (errno int) = SYS_FSTAT64
//sys	Fstatfs(fd int, stat *Statfs_t) (errno int) = SYS_FSTATFS64
//sys	Fsync(fd int) (errno int)
//sys	Ftruncate(fd int, length int64) (errno int)
//sys	Getdirentries(fd int, buf []byte, basep *uintptr) (n int, errno int) = SYS_GETDIRENTRIES64
//sys	Getdtablesize() (size int)
//sys	Getegid() (egid int)
//sys	Geteuid() (uid int)
//sys	Getfsstat(buf []Statfs_t, flags int) (n int, errno int) = SYS_GETFSSTAT64
//sys	Getgid() (gid int)
//sys	Getpgid(pid int) (pgid int, errno int)
//sys	Getpgrp() (pgrp int)
//sys	Getpid() (pid int)
//sys	Getppid() (ppid int)
//sys	Getpriority(which int, who int) (prio int, errno int)
//sys	Getrlimit(which int, lim *Rlimit) (errno int)
//sys	Getrusage(who int, rusage *Rusage) (errno int)
//sys	Getsid(pid int) (sid int, errno int)
//sys	Getuid() (uid int)
//sys	Issetugid() (tainted bool)
//sys	Kill(pid int, signum int, posix int) (errno int)
//sys	Kqueue() (fd int, errno int)
//sys	Lchown(path string, uid int, gid int) (errno int)
//sys	Link(path string, link string) (errno int)
//sys	Listen(s int, backlog int) (errno int)
//sys	Lstat(path string, stat *Stat_t) (errno int) = SYS_LSTAT64
//sys	Mkdir(path string, mode int) (errno int)
//sys	Mkfifo(path string, mode int) (errno int)
//sys	Mknod(path string, mode int, dev int) (errno int)
//sys	Open(path string, mode int, perm int) (fd int, errno int)
//sys	Pathconf(path string, name int) (val int, errno int)
//sys	Pread(fd int, p []byte, offset int64) (n int, errno int)
//sys	Pwrite(fd int, p []byte, offset int64) (n int, errno int)
//sys	Read(fd int, p []byte) (n int, errno int)
//sys	Readlink(path string, buf []byte) (n int, errno int)
//sys	Rename(from string, to string) (errno int)
//sys	Revoke(path string) (errno int)
//sys	Rmdir(path string) (errno int)
//sys	Seek(fd int, offset int64, whence int) (newoffset int64, errno int) = SYS_LSEEK
//sys	Select(n int, r *FdSet, w *FdSet, e *FdSet, timeout *Timeval) (errno int)
//sys	Setegid(egid int) (errno int)
//sys	Seteuid(euid int) (errno int)
//sys	Setgid(gid int) (errno int)
//sys	Setlogin(name string) (errno int)
//sys	Setpgid(pid int, pgid int) (errno int)
//sys	Setpriority(which int, who int, prio int) (errno int)
//sys	Setprivexec(flag int) (errno int)
//sys	Setregid(rgid int, egid int) (errno int)
//sys	Setreuid(ruid int, euid int) (errno int)
//sys	Setrlimit(which int, lim *Rlimit) (errno int)
//sys	Setsid() (pid int, errno int)
//sys	Settimeofday(tp *Timeval) (errno int)
//sys	Setuid(uid int) (errno int)
//sys	Stat(path string, stat *Stat_t) (errno int) = SYS_STAT64
//sys	Statfs(path string, stat *Statfs_t) (errno int) = SYS_STATFS64
//sys	Symlink(path string, link string) (errno int)
//sys	Sync() (errno int)
//sys	Truncate(path string, length int64) (errno int)
//sys	Umask(newmask int) (errno int)
//sys	Undelete(path string) (errno int)
//sys	Unlink(path string) (errno int)
//sys	Unmount(path string, flags int) (errno int)
//sys	Write(fd int, p []byte) (n int, errno int)
//sys	read(fd int, buf *byte, nbuf int) (n int, errno int)
//sys	write(fd int, buf *byte, nbuf int) (n int, errno int)


/*
 * Unimplemented
 */
// Profil
// Sigaction
// Sigprocmask
// Getlogin
// Sigpending
// Sigaltstack
// Ioctl
// Reboot
// Execve
// Vfork
// Sbrk
// Sstk
// Ovadvise
// Mincore
// Setitimer
// Swapon
// Select
// Sigsuspend
// Readv
// Writev
// Nfssvc
// Getfh
// Quotactl
// Mount
// Csops
// Waitid
// Add_profil
// Kdebug_trace
// Sigreturn
// Mmap
// __Sysctl
// Mlock
// Munlock
// Atsocket
// Kqueue_from_portset_np
// Kqueue_portset
// Getattrlist
// Setattrlist
// Getdirentriesattr
// Searchfs
// Delete
// Copyfile
// Poll
// Watchevent
// Waitevent
// Modwatch
// Getxattr
// Fgetxattr
// Setxattr
// Fsetxattr
// Removexattr
// Fremovexattr
// Listxattr
// Flistxattr
// Fsctl
// Initgroups
// Posix_spawn
// Nfsclnt
// Fhopen
// Minherit
// Semsys
// Msgsys
// Shmsys
// Semctl
// Semget
// Semop
// Msgctl
// Msgget
// Msgsnd
// Msgrcv
// Shmat
// Shmctl
// Shmdt
// Shmget
// Shm_open
// Shm_unlink
// Sem_open
// Sem_close
// Sem_unlink
// Sem_wait
// Sem_trywait
// Sem_post
// Sem_getvalue
// Sem_init
// Sem_destroy
// Open_extended
// Umask_extended
// Stat_extended
// Lstat_extended
// Fstat_extended
// Chmod_extended
// Fchmod_extended
// Access_extended
// Settid
// Gettid
// Setsgroups
// Getsgroups
// Setwgroups
// Getwgroups
// Mkfifo_extended
// Mkdir_extended
// Identitysvc
// Shared_region_check_np
// Shared_region_map_np
// __pthread_mutex_destroy
// __pthread_mutex_init
// __pthread_mutex_lock
// __pthread_mutex_trylock
// __pthread_mutex_unlock
// __pthread_cond_init
// __pthread_cond_destroy
// __pthread_cond_broadcast
// __pthread_cond_signal
// Setsid_with_pid
// __pthread_cond_timedwait
// Aio_fsync
// Aio_return
// Aio_suspend
// Aio_cancel
// Aio_error
// Aio_read
// Aio_write
// Lio_listio
// __pthread_cond_wait
// Iopolicysys
// Mlockall
// Munlockall
// __pthread_kill
// __pthread_sigmask
// __sigwait
// __disable_threadsignal
// __pthread_markcancel
// __pthread_canceled
// __semwait_signal
// Proc_info
// Sendfile
// Stat64_extended
// Lstat64_extended
// Fstat64_extended
// __pthread_chdir
// __pthread_fchdir
// Audit
// Auditon
// Getauid
// Setauid
// Getaudit
// Setaudit
// Getaudit_addr
// Setaudit_addr
// Auditctl
// Bsdthread_create
// Bsdthread_terminate
// Stack_snapshot
// Bsdthread_register
// Workq_open
// Workq_ops
// __mac_execve
// __mac_syscall
// __mac_get_file
// __mac_set_file
// __mac_get_link
// __mac_set_link
// __mac_get_proc
// __mac_set_proc
// __mac_get_fd
// __mac_set_fd
// __mac_get_pid
// __mac_get_lcid
// __mac_get_lctx
// __mac_set_lctx
// Setlcid
// Read_nocancel
// Write_nocancel
// Open_nocancel
// Close_nocancel
// Wait4_nocancel
// Recvmsg_nocancel
// Sendmsg_nocancel
// Recvfrom_nocancel
// Accept_nocancel
// Msync_nocancel
// Fcntl_nocancel
// Select_nocancel
// Fsync_nocancel
// Connect_nocancel
// Sigsuspend_nocancel
// Readv_nocancel
// Writev_nocancel
// Sendto_nocancel
// Pread_nocancel
// Pwrite_nocancel
// Waitid_nocancel
// Poll_nocancel
// Msgsnd_nocancel
// Msgrcv_nocancel
// Sem_wait_nocancel
// Aio_suspend_nocancel
// __sigwait_nocancel
// __semwait_signal_nocancel
// __mac_mount
// __mac_get_mount
// __mac_getfsstat

