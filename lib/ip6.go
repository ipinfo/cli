package lib

import (
	"encoding/binary"
	"net"
)

// IP6 is a 128-bit number representation of a IPv6 address in big endian byte
// order.
//
// The number is internally represented as 2 64-bit numbers.
type IP6 struct {
	N U128
}

// NewIP6 returns a new IP6 based off of `hi` and `lo` numeric parts.
func NewIP6(hi uint64, lo uint64) IP6 {
	return IP6{N: U128{Hi: hi, Lo: lo}}
}

// IP6FromStdIP returns a new IP6 from a standard library `net.IP`, and whether
// the conversion succeeded.
func IP6FromStdIP(ip net.IP) (IP6, bool) {
	switch len(ip) {
	case 16:
		return IP6FromByteSlice(ip), true
	case 4:
		return IP6FromIP4Bytes(ip[0], ip[1], ip[2], ip[3]), true
	}
	return IP6{}, false
}

// IP6FromBytes returns a new IP6 given 16 bytes representing the 16 bytes of
// the IPv6 address in big-endian order.
func IP6FromBytes(b [16]byte) IP6 {
	return IP6{N: U128{
		Hi: binary.BigEndian.Uint64(b[:8]),
		Lo: binary.BigEndian.Uint64(b[8:]),
	}}
}

// IP6FromByteSlice is the same as IP6FromBytes but from a byte slice which may
// be longer than 16 bytes long.
//
// <16 byte slice will cause a panic.
func IP6FromByteSlice(b []byte) IP6 {
	return IP6{N: U128{
		Hi: binary.BigEndian.Uint64(b[:8]),
		Lo: binary.BigEndian.Uint64(b[8:]),
	}}
}

// IP6FromIP4Bytes returns a new IP6 from the IPv4 address of bytes `a.b.c.d`.
func IP6FromIP4Bytes(a uint8, b uint8, c uint8, d uint8) IP6 {
	return IP6{N: U128{
		Hi: 0,
		Lo: 0xffff00000000 | uint64(a)<<24 | uint64(b)<<16 | uint64(c)<<8 | uint64(d),
	}}
}

// To16Bytes returns the 16-byte array representation of an IPv6 address in
// big-endian byte order.
func (ip IP6) To16Bytes() [16]byte {
	var b [16]byte
	binary.BigEndian.PutUint64(b[:8], ip.N.Hi)
	binary.BigEndian.PutUint64(b[8:], ip.N.Lo)
	return b
}

// To16ByteSlice returns the 16-byte slice representation of an IPv6 address in
// big-endian byte order.
func (ip IP6) To16ByteSlice() []byte {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[:8], ip.N.Hi)
	binary.BigEndian.PutUint64(b[8:], ip.N.Lo)
	return b
}

// ToStdIP returns the 16-byte slice standard library IPv6 representation.
func (ip IP6) ToStdIP() net.IP {
	return net.IP(ip.To16ByteSlice())
}

// Cmp compares `ip1` and `ip2` and returns:
//
// - -1 if `ip1<ip2`
// -  0 if `ip1==ip2`
// -  1 if `ip1>ip2`
func (ip1 IP6) Cmp(ip2 IP6) int {
	return ip1.N.Cmp(ip2.N)
}

// Eq returns whether `ip1 == ip2` numerically.
func (ip1 IP6) Eq(ip2 IP6) bool {
	return ip1.N.Eq(ip2.N)
}

// Gt returns whether `ip1 > ip2` numerically.
func (ip1 IP6) Gt(ip2 IP6) bool {
	return ip1.N.Gt(ip2.N)
}

// Gte returns whether `ip1 >= ip2` numerically.
func (ip1 IP6) Gte(ip2 IP6) bool {
	return ip1.N.Gte(ip2.N)
}

// Lt returns whether `ip1 < ip2` numerically.
func (ip1 IP6) Lt(ip2 IP6) bool {
	return ip1.N.Lt(ip2.N)
}

// Lte returns whether `ip1 <= ip2` numerically.
func (ip1 IP6) Lte(ip2 IP6) bool {
	return ip1.N.Lte(ip2.N)
}

// String returns the IPv6 string representation of `ip`.
func (ip IP6) String() string {
	return ip.ToStdIP().String()
}
