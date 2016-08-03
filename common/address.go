package stun

import (
	"fmt"
	"net"
)

type Address struct {
	net.Addr
}

func (a *Address) HostPart() net.Addr {
	host, _, err := net.SplitHostPort(a.String())
	if err != nil {
		return a
	} else {
		addr, _ := net.ResolveIPAddr("ip", host)
		return addr
	}
}

// The STUN family of an address. A value of 1 if IPv4 and 2 if IPv6
func (a *Address) Family() uint16 {
	switch a.Network() {
	case "ip6", "udp6", "tcp6":
		return uint16(2)
	default:
		return uint16(1)
	}
}

func (a *Address) Port() uint16 {
	switch a.Network() {
	case "tcp", "tcp4", "tcp6":
		addr, _ := net.ResolveTCPAddr(a.Network(), a.String())
		return uint16(addr.Port)
	case "udp", "udp4", "udp6":
		addr, _ := net.ResolveUDPAddr(a.Network(), a.String())
		return uint16(addr.Port)
	default:
		return uint16(0)
	}
}

func (a *Address) Host() net.IP {
	switch a.Network() {
	case "tcp", "tcp4", "tcp6":
		addr, _ := net.ResolveTCPAddr(a.Network(), a.String())
		return addr.IP
	case "udp", "udp4", "udp6":
		addr, _ := net.ResolveUDPAddr(a.Network(), a.String())
		return addr.IP
	case "ip", "ip4", "ip6":
		addr, _ := net.ResolveIPAddr(a.Network(), a.String())
		return addr.IP
	default:
		return net.IP{}
	}
}

func NewAddress(network string, host net.IP, port uint16) Address {
	hostport := net.JoinHostPort(host.String(), fmt.Sprintf("%d", port))
	switch network {
	case "tcp", "tcp4", "tcp6":
		addr, _ := net.ResolveTCPAddr(network, hostport)
		return Address{addr}
	case "udp", "udp4", "udp6":
		addr, _ := net.ResolveUDPAddr(network, hostport)
		return Address{addr}
	default:
		return Address{}
	}
}
