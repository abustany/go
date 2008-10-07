// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"os";
	"net";
	"strings";
	"syscall"
)

export var (
	BadAddress = os.NewError("malformed address");
	MissingAddress = os.NewError("missing address");
	UnknownNetwork = os.NewError("unknown network");
	UnknownHost = os.NewError("unknown host");
	UnknownPort = os.NewError("unknown port");
	UnknownSocketFamily = os.NewError("unknown socket family");
)

// Split "host:port" into "host" and "port".
// Host cannot contain colons unless it is bracketed.
func SplitHostPort(hostport string) (host, port string, err *os.Error) {
	// The port starts after the last colon.
	var i int;
	for i = len(hostport)-1; i >= 0; i-- {
		if hostport[i] == ':' {
			break
		}
	}
	if i < 0 {
		return "", "", BadAddress
	}

	host = hostport[0:i];
	port = hostport[i+1:len(hostport)];

	// Can put brackets around host ...
	if host[0] == '[' && host[len(host)-1] == ']' {
		host = host[1:len(host)-1]
	} else {
		// ... but if there are no brackets, no colons.
		for i := 0; i < len(host); i++ {
			if host[i] == ':' {
				return "", "", BadAddress
			}
		}
	}
	return host, port, nil
}

// Join "host" and "port" into "host:port".
// If host contains colons, will join into "[host]:port".
func JoinHostPort(host, port string) string {
	// If host has colons, have to bracket it.
	for i := 0; i < len(host); i++ {
		if host[i] == ':' {
			return "[" + host + "]:" + port
		}
	}
	return host + ":" + port
}

func xdtoi(s string) (n int, ok bool) {
	if s == "" || s[0] < '0' || s[0] > '9' {
		return 0, false
	}
	n = 0;
	for i := 0; i < len(s) && '0' <= s[i] && s[i] <= '9'; i++ {
		n = n*10 + int(s[i] - '0');
		if n >= 1000000 {	// bigger than we need
			return 0, false
		}
	}
	return n, true
}

// Convert "host:port" into IP address and port.
// For now, host and port must be numeric literals.
// Eventually, we'll have name resolution.
func HostPortToIP(net string, hostport string) (ip *[]byte, iport int, err *os.Error) {
	var host, port string;
	host, port, err = SplitHostPort(hostport);
	if err != nil {
		return nil, 0, err
	}

	// TODO: Resolve host.

	addr := ParseIP(host);
	if addr == nil {
		return nil, 0, UnknownHost
	}

	// TODO: Resolve port.

	p, ok := xdtoi(port);
	if !ok || p < 0 || p > 0xFFFF {
		return nil, 0, UnknownPort
	}

	return addr, p, nil
}

// Convert socket address into "host:port".
func SockaddrToHostPort(sa *syscall.Sockaddr) (hostport string, err *os.Error) {
	switch sa.family {
	case syscall.AF_INET, syscall.AF_INET6:
		addr, port, e := SockaddrToIP(sa);
		if e != nil {
			return "", e
		}
		host := IPToString(addr);
		return JoinHostPort(host, strings.itoa(port)), nil;
	default:
		return "", UnknownSocketFamily
	}
	return "", nil // not reached
}

// Boolean to int.
func boolint(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Generic Socket creation.
func Socket(f, p, t int64, la, ra *syscall.Sockaddr) (fd *FD, err *os.Error) {
	s, e := syscall.socket(f, p, t);
	if e != 0 {
		return nil, os.ErrnoToError(e)
	}

	// Allow reuse of recently-used addresses.
	syscall.setsockopt_int(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1);

	var r int64;
	if la != nil {
		r, e = syscall.bind(s, la);
		if e != 0 {
			syscall.close(s);
			return nil, os.ErrnoToError(e)
		}
	}

	if ra != nil {
		r, e = syscall.connect(s, ra);
		if e != 0 {
			syscall.close(s);
			return nil, os.ErrnoToError(e)
		}
	}

	fd, err = NewFD(s);
	if err != nil {
		syscall.close(s);
		return nil, err
	}

	return fd, nil
}


// Generic implementation of Conn interface; not exported.
type ConnBase struct {
	fd *FD;
	raddr string;
}

func (c *ConnBase) FD() int64 {
	if c == nil || c.fd == nil {
		return -1
	}
	return c.fd.fd
}

func (c *ConnBase) Read(b *[]byte) (n int, err *os.Error) {
	n, err = c.fd.Read(b);
	return n, err
}

func (c *ConnBase) Write(b *[]byte) (n int, err *os.Error) {
	n, err = c.fd.Write(b);
	return n, err
}

func (c *ConnBase) ReadFrom(b *[]byte) (n int, raddr string, err *os.Error) {
	if c == nil {
		return -1, "", os.EINVAL
	}
	n, err = c.Read(b);
	return n, c.raddr, err
}

func (c *ConnBase) WriteTo(raddr string, b *[]byte) (n int, err *os.Error) {
	if c == nil {
		return -1, os.EINVAL
	}
	if raddr != c.raddr {
		return -1, os.EINVAL
	}
	n, err = c.Write(b);
	return n, err
}

func (c *ConnBase) Close() *os.Error {
	if c == nil {
		return os.EINVAL
	}
	return c.fd.Close()
}


func setsockopt_int(fd, level, opt int64, value int) *os.Error {
	return os.ErrnoToError(syscall.setsockopt_int(fd, level, opt, value));
}

func setsockopt_tv(fd, level, opt int64, nsec int64) *os.Error {
	return os.ErrnoToError(syscall.setsockopt_tv(fd, level, opt, nsec));
}

func (c *ConnBase) SetReadBuffer(bytes int) *os.Error {
	return setsockopt_int(c.FD(), syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes);
}

func (c *ConnBase) SetWriteBuffer(bytes int) *os.Error {
	return setsockopt_int(c.FD(), syscall.SOL_SOCKET, syscall.SO_SNDBUF, bytes);
}

func (c *ConnBase) SetReadTimeout(nsec int64) *os.Error {
	return setsockopt_tv(c.FD(), syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, nsec);
}

func (c *ConnBase) SetWriteTimeout(nsec int64) *os.Error {
	return setsockopt_tv(c.FD(), syscall.SOL_SOCKET, syscall.SO_SNDTIMEO, nsec);
}

func (c *ConnBase) SetTimeout(nsec int64) *os.Error {
	if e := c.SetReadTimeout(nsec); e != nil {
		return e
	}
	return c.SetWriteTimeout(nsec)
}

func (c *ConnBase) SetReuseAddr(reuse bool) *os.Error {
	return setsockopt_int(c.FD(), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, boolint(reuse));
}

func (c *ConnBase) BindToDevice(dev string) *os.Error {
	// TODO: call setsockopt with null-terminated string pointer
	return os.EINVAL
}

func (c *ConnBase) SetDontRoute(dontroute bool) *os.Error {
	return setsockopt_int(c.FD(), syscall.SOL_SOCKET, syscall.SO_DONTROUTE, boolint(dontroute));
}

func (c *ConnBase) SetKeepAlive(keepalive bool) *os.Error {
	return setsockopt_int(c.FD(), syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, boolint(keepalive));
}

func (c *ConnBase) SetLinger(sec int) *os.Error {
	e := syscall.setsockopt_linger(c.FD(), syscall.SOL_SOCKET, syscall.SO_LINGER, sec);
	return os.ErrnoToError(e);
}


// Internet sockets (TCP, UDP)

// Should we try to use the IPv4 socket interface if we're
// only dealing with IPv4 sockets?  As long as the host system
// understands IPv6, it's okay to pass IPv4 addresses to the IPv6
// interface.  That simplifies our code and is most general.
// If we need to build on a system without IPv6 support, setting
// PreferIPv4 here should fall back to the IPv4 socket interface when possible.
const PreferIPv4 = false

func InternetSocket(net, laddr, raddr string, proto int64) (fd *FD, err *os.Error) {
	// Parse addresses (unless they are empty).
	var lip, rip *[]byte;
	var lport, rport int;
	var lerr, rerr *os.Error;
// BUG 6g doesn't zero var lists
lip = nil;
rip = nil;
lport = 0;
rport = 0;
lerr = nil;
rerr = nil;
	if laddr != "" {
		lip, lport, lerr = HostPortToIP(net, laddr);
		if lerr != nil {
			return nil, lerr
		}
	}
	if raddr != "" {
		rip, rport, rerr = HostPortToIP(net, raddr);
		if rerr != nil {
			return nil, rerr
		}
	}

	// Figure out IP version.
	// If network has a suffix like "tcp4", obey it.
	vers := 0;
	switch net[len(net)-1] {
	case '4':
		vers = 4;
	case '6':
		vers = 6;
	default:
		// Otherwise, guess.
		// If the addresses are IPv4 and we prefer IPv4, use 4; else 6.
		if PreferIPv4
		&& (lip == nil || ToIPv4(lip) != nil)
		&& (rip == nil || ToIPv4(rip) != nil) {
			vers = 4
		} else {
			vers = 6
		}
	}

	var cvt *(addr *[]byte, port int) (sa *syscall.Sockaddr, err *os.Error);
	var family int64;
	if vers == 4 {
		cvt = &IPv4ToSockaddr;
		family = syscall.AF_INET
	} else {
		cvt = &IPv6ToSockaddr;
		family = syscall.AF_INET6
	}

	var la, ra *syscall.Sockaddr;
// BUG
la = nil;
ra = nil;
	if lip != nil {
		la, lerr = cvt(lip, lport);
		if lerr != nil {
			return nil, lerr
		}
	}
	if rip != nil {
		ra, rerr = cvt(rip, rport);
		if rerr != nil {
			return nil, rerr
		}
	}

	fd, err = Socket(family, proto, 0, la, ra);
	return fd, err
}


// TCP connections.

export type ConnTCP struct {
	base ConnBase
}

// New TCP methods
func (c *ConnTCP) SetNoDelay(nodelay bool) *os.Error {
	if c == nil {
		return os.EINVAL
	}
	return setsockopt_int((&c.base).FD(), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, boolint(nodelay))
}

// Wrappers
func (c *ConnTCP) Read(b *[]byte) (n int, err *os.Error) {
	n, err = (&c.base).Read(b);
	return n, err
}
func (c *ConnTCP) Write(b *[]byte) (n int, err *os.Error) {
	n, err = (&c.base).Write(b);
	return n, err
}
func (c *ConnTCP) ReadFrom(b *[]byte) (n int, raddr string, err *os.Error) {
	n, raddr, err = (&c.base).ReadFrom(b);
	return n, raddr, err
}
func (c *ConnTCP) WriteTo(raddr string, b *[]byte) (n int, err *os.Error) {
	n, err = (&c.base).WriteTo(raddr, b);
	return n, err
}
func (c *ConnTCP) Close() *os.Error {
	return (&c.base).Close()
}
func (c *ConnTCP) SetReadBuffer(bytes int) *os.Error {
	return (&c.base).SetReadBuffer(bytes)
}
func (c *ConnTCP) SetWriteBuffer(bytes int) *os.Error {
	return (&c.base).SetWriteBuffer(bytes)
}
func (c *ConnTCP) SetTimeout(nsec int64) *os.Error {
	return (&c.base).SetTimeout(nsec)
}
func (c *ConnTCP) SetReadTimeout(nsec int64) *os.Error {
	return (&c.base).SetReadTimeout(nsec)
}
func (c *ConnTCP) SetWriteTimeout(nsec int64) *os.Error {
	return (&c.base).SetWriteTimeout(nsec)
}
func (c *ConnTCP) SetLinger(sec int) *os.Error {
	return (&c.base).SetLinger(sec)
}
func (c *ConnTCP) SetReuseAddr(reuseaddr bool) *os.Error {
	return (&c.base).SetReuseAddr(reuseaddr)
}
func (c *ConnTCP) BindToDevice(dev string) *os.Error {
	return (&c.base).BindToDevice(dev)
}
func (c *ConnTCP) SetDontRoute(dontroute bool) *os.Error {
	return (&c.base).SetDontRoute(dontroute)
}
func (c *ConnTCP) SetKeepAlive(keepalive bool) *os.Error {
	return (&c.base).SetKeepAlive(keepalive)
}

func NewConnTCP(fd *FD, raddr string) *ConnTCP {
	c := new(ConnTCP);
	c.base.fd = fd;
	c.base.raddr = raddr;
	c.SetNoDelay(true);
	return c
}

export func DialTCP(net, laddr, raddr string) (c *ConnTCP, err *os.Error) {
	if raddr == "" {
		return nil, MissingAddress
	}
	fd, e := InternetSocket(net, laddr, raddr, syscall.SOCK_STREAM);
	if e != nil {
		return nil, e
	}
	return NewConnTCP(fd, raddr), nil
}


// TODO: UDP connections


// TODO: raw IP connections


// TODO: raw ethernet connections


export type Conn interface {
	Read(b *[]byte) (n int, err *os.Error);
	Write(b *[]byte) (n int, err *os.Error);
	ReadFrom(b *[]byte) (n int, addr string, err *os.Error);
	WriteTo(addr string, b *[]byte) (n int, err *os.Error);
	Close() *os.Error;
	SetReadBuffer(bytes int) *os.Error;
	SetWriteBuffer(bytes int) *os.Error;
	SetTimeout(nsec int64) *os.Error;
	SetReadTimeout(nsec int64) *os.Error;
	SetWriteTimeout(nsec int64) *os.Error;
	SetLinger(sec int) *os.Error;
	SetReuseAddr(reuseaddr bool) *os.Error;
	SetDontRoute(dontroute bool) *os.Error;
	SetKeepAlive(keepalive bool) *os.Error;
	BindToDevice(dev string) *os.Error;
}

type NoConn struct { unused int }
func (c *NoConn) Read(b *[]byte) (n int, err *os.Error) { return -1, os.EINVAL }
func (c *NoConn) Write(b *[]byte) (n int, err *os.Error) { return -1, os.EINVAL }
func (c *NoConn) ReadFrom(b *[]byte) (n int, addr string, err *os.Error) { return -1, "", os.EINVAL }
func (c *NoConn) WriteTo(addr string, b *[]byte) (n int, err *os.Error) { return -1, os.EINVAL }
func (c *NoConn) Close() *os.Error { return nil }
func (c *NoConn) SetReadBuffer(bytes int) *os.Error { return os.EINVAL }
func (c *NoConn) SetWriteBuffer(bytes int) *os.Error { return os.EINVAL }
func (c *NoConn) SetTimeout(nsec int64) *os.Error { return os.EINVAL }
func (c *NoConn) SetReadTimeout(nsec int64) *os.Error { return os.EINVAL }
func (c *NoConn) SetWriteTimeout(nsec int64) *os.Error { return os.EINVAL }
func (c *NoConn) SetLinger(sec int) *os.Error { return os.EINVAL }
func (c *NoConn) SetReuseAddr(reuseaddr bool) *os.Error { return os.EINVAL }
func (c *NoConn) SetDontRoute(dontroute bool) *os.Error { return os.EINVAL }
func (c *NoConn) SetKeepAlive(keepalive bool) *os.Error { return os.EINVAL }
func (c *NoConn) BindToDevice(dev string) *os.Error { return os.EINVAL }

var noconn NoConn

// Dial's arguments are the network, local address, and remote address.
// Examples:
//	Dial("tcp", "", "12.34.56.78:80")
//	Dial("tcp", "", "[de:ad:be:ef::ca:fe]:80")
//	Dial("tcp", "127.0.0.1:123", "127.0.0.1:88")
//
// Eventually, we plan to allow names in addition to IP addresses,
// but that requires writing a DNS library.

export func Dial(net, laddr, raddr string) (c Conn, err *os.Error) {
	switch net {
	case "tcp", "tcp4", "tcp6":
		c, err := DialTCP(net, laddr, raddr);
		if err != nil {
			return &noconn, err
		}
		return c, nil;
/*
	case "udp", "udp4", "upd6":
		c, err := DialUDP(net, laddr, raddr);
		return c, err;
	case "ether":
		c, err := DialEther(net, laddr, raddr);
		return c, err;
	case "ipv4":
		c, err := DialIPv4(net, laddr, raddr);
		return c, err;
	case "ipv6":
		c, err := DialIPv6(net, laddr, raddr);
		return c, err
*/
	}
	return nil, UnknownNetwork
}


export type Listener interface {
	Accept() (c Conn, raddr string, err *os.Error);
	Close() *os.Error;
}

type NoListener struct { unused int }
func (l *NoListener) Accept() (c Conn, raddr string, err *os.Error) {
	return &noconn, "", os.EINVAL
}
func (l *NoListener) Close() *os.Error { return os.EINVAL }

var nolistener NoListener

export type ListenerTCP struct {
	fd *FD;
	laddr string
}

export func ListenTCP(net, laddr string) (l *ListenerTCP, err *os.Error) {
	fd, e := InternetSocket(net, laddr, "", syscall.SOCK_STREAM);
	if e != nil {
		return nil, e
	}
	r, e1 := syscall.listen(fd.fd, ListenBacklog());
	if e1 != 0 {
		syscall.close(fd.fd);
		return nil, os.ErrnoToError(e1)
	}
	l = new(ListenerTCP);
	l.fd = fd;
	return l, nil
}

func (l *ListenerTCP) AcceptTCP() (c *ConnTCP, raddr string, err *os.Error) {
	if l == nil || l.fd == nil || l.fd.fd < 0 {
		return nil, "", os.EINVAL
	}
	var sa syscall.Sockaddr;
	fd, e := l.fd.Accept(&sa);
	if e != nil {
		return nil, "", e
	}
	raddr, err = SockaddrToHostPort(&sa);
	if err != nil {
		fd.Close();
		return nil, "", err
	}
	return NewConnTCP(fd, raddr), raddr, nil
}

func (l *ListenerTCP) Accept() (c Conn, raddr string, err *os.Error) {
	c1, r1, e1 := l.AcceptTCP();
	if e1 != nil {
		return &noconn, "", e1
	}
	return c1, r1, nil
}

func (l *ListenerTCP) Close() *os.Error {
	if l == nil || l.fd == nil {
		return os.EINVAL
	}
	return l.fd.Close()
}

export func Listen(net, laddr string) (l Listener, err *os.Error) {
	switch net {
	case "tcp", "tcp4", "tcp6":
		l, err := ListenTCP(net, laddr);
		if err != nil {
			return &nolistener, err
		}
		return l, nil
/*
	more here
*/
	}
	return nil, UnknownNetwork
}

