package lib

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
)

// IP is a numerical representation of an IPv4 address.
// The number must be in big-endian form.
type IP uint32

// IPStr is a string representation of an IPv4 address.
type IPStr string

// NewIP returns a new IPv4 address representation.
// `ip` must already be in big-endian form.
func NewIP(ip uint32) IP {
	return IP(ip)
}

// RandIP4 returns a new randomly generated IPv4 address.
func RandIP4() net.IP {
	ip := make([]byte, 4)
	binary.BigEndian.PutUint32(ip, rand.Uint32())
	return net.IPv4(ip[0], ip[1], ip[2], ip[3])
}

// RandIP4List returns a list of new randomly generated IPv4 addresses.
func RandIP4List(n int) []net.IP {
	ip := make([]net.IP, n)
	for i := 0; i < n; i++ {
		ip[i] = RandIP4()
	}
	return ip
}

// RandIP4Write prints IP/list of new randomly generated IPv4 addresses.
func RandIP4Write(n int) {
	for i := 0; i < n; i++ {
		fmt.Println(RandIP4())
	}
}

// IPFromStdIP returns a new IPv4 address representation from the standard
// library's IP representation.
func IPFromStdIP(ip net.IP) IP {
	return IP(binary.BigEndian.Uint32(ip.To4()))
}

// String returns the IPv4 string representation of `ip`.
func (ip IP) String() string {
	return net.IPv4(
		byte(ip>>24),
		byte(ip>>16),
		byte(ip>>8),
		byte(ip),
	).String()
}
