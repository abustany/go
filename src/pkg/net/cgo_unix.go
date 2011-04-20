// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

/*
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
*/
import "C"

import (
	"os"
	"syscall"
	"unsafe"
)

func cgoLookupHost(name string) (addrs []string, err os.Error, completed bool) {
	ip, err, completed := cgoLookupIP(name)
	for _, p := range ip {
		addrs = append(addrs, p.String())
	}
	return
}

func cgoLookupPort(net, service string) (port int, err os.Error, completed bool) {
	var res *C.struct_addrinfo
	var hints C.struct_addrinfo

	switch net {
	case "":
		// no hints
	case "tcp", "tcp4", "tcp6":
		hints.ai_socktype = C.SOCK_STREAM
		hints.ai_protocol = C.IPPROTO_TCP
	case "udp", "udp4", "udp6":
		hints.ai_socktype = C.SOCK_DGRAM
		hints.ai_protocol = C.IPPROTO_UDP
	default:
		return 0, UnknownNetworkError(net), true
	}
	if len(net) >= 4 {
		switch net[3] {
		case '4':
			hints.ai_family = C.AF_INET
		case '6':
			hints.ai_family = C.AF_INET6
		}
	}

	s := C.CString(service)
	defer C.free(unsafe.Pointer(s))
	if C.getaddrinfo(nil, s, &hints, &res) == 0 {
		defer C.freeaddrinfo(res)
		for r := res; r != nil; r = r.ai_next {
			switch r.ai_family {
			default:
				continue
			case C.AF_INET:
				sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(r.ai_addr))
				p := (*[2]byte)(unsafe.Pointer(&sa.Port))
				return int(p[0])<<8 | int(p[1]), nil, true
			case C.AF_INET6:
				sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(r.ai_addr))
				p := (*[2]byte)(unsafe.Pointer(&sa.Port))
				return int(p[0])<<8 | int(p[1]), nil, true
			}
		}
	}
	return 0, &AddrError{"unknown port", net + "/" + service}, true
}

func cgoLookupIPCNAME(name string) (addrs []IP, cname string, err os.Error, completed bool) {
	var res *C.struct_addrinfo
	var hints C.struct_addrinfo

	// NOTE(rsc): In theory there are approximately balanced
	// arguments for and against including AI_ADDRCONFIG
	// in the flags (it includes IPv4 results only on IPv4 systems,
	// and similarly for IPv6), but in practice setting it causes
	// getaddrinfo to return the wrong canonical name on Linux.
	// So definitely leave it out.
	hints.ai_flags = C.AI_ALL | C.AI_V4MAPPED | C.AI_CANONNAME

	h := C.CString(name)
	defer C.free(unsafe.Pointer(h))
	gerrno, err := C.getaddrinfo(h, nil, &hints, &res)
	if gerrno != 0 {
		var str string
		if gerrno == C.EAI_NONAME {
			str = noSuchHost
		} else if gerrno == C.EAI_SYSTEM {
			str = err.String()
		} else {
			str = C.GoString(C.gai_strerror(gerrno))
		}
		return nil, "", &DNSError{Error: str, Name: name}, true
	}
	defer C.freeaddrinfo(res)
	if res != nil {
		cname = C.GoString(res.ai_canonname)
		if cname == "" {
			cname = name
		}
		if len(cname) > 0 && cname[len(cname)-1] != '.' {
			cname += "."
		}
	}
	for r := res; r != nil; r = r.ai_next {
		// Everything comes back twice, once for UDP and once for TCP.
		if r.ai_socktype != C.SOCK_STREAM {
			continue
		}
		switch r.ai_family {
		default:
			continue
		case C.AF_INET:
			sa := (*syscall.RawSockaddrInet4)(unsafe.Pointer(r.ai_addr))
			addrs = append(addrs, copyIP(sa.Addr[:]))
		case C.AF_INET6:
			sa := (*syscall.RawSockaddrInet6)(unsafe.Pointer(r.ai_addr))
			addrs = append(addrs, copyIP(sa.Addr[:]))
		}
	}
	return addrs, cname, nil, true
}

func cgoLookupIP(name string) (addrs []IP, err os.Error, completed bool) {
	addrs, _, err, completed = cgoLookupIPCNAME(name)
	return
}

func cgoLookupCNAME(name string) (cname string, err os.Error, completed bool) {
	_, cname, err, completed = cgoLookupIPCNAME(name)
	return
}

func copyIP(x IP) IP {
	y := make(IP, len(x))
	copy(y, x)
	return y
}
