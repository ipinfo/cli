package lib

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net"
)

// IP is a numerical representation of an IPv4 address.
// The number must be in big-endian form.
type IP uint32

// IPStr is a string representation of an IPv4 address.
type IPStr string

// IP4Range is a numerical representation of an IPv4 range.
type IP4Range struct {
	startIP IP
	endIP   IP
}

// IP6Range is a numerical representation of an IPv6 range.
type IP6RangeInt struct {
	startIP *big.Int
	endIP   *big.Int
}

// NewIP returns a new IPv4 address representation.
// `ip` must already be in big-endian form.
func NewIP(ip uint32) IP {
	return IP(ip)
}

// RandIP4 returns a new randomly generated IPv4 address.
func RandIP4() net.IP {
	ip := [4]byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(ip[:], rand.Uint32())
	return net.IPv4(ip[0], ip[1], ip[2], ip[3])
}

// RandIP4List returns a list of new randomly generated IPv4 addresses.
func RandIP4List(n int) []net.IP {
	ips := make([]net.IP, n)
	for i := 0; i < n; i++ {
		ips[i] = RandIP4()
	}
	return ips
}

// NewIP4Range returns a new IP4Range addresses representation.
//
// note: starting and ending IPs must be valid IPv4 string formats.
func NewIP4Range(
	startIP string,
	endIP string,
) (IP4Range, error) {
	startIPRaw := net.ParseIP(startIP).To4()
	if len(startIPRaw) == 0 || len(startIPRaw) > net.IPv4len {
		return IP4Range{}, errors.New("invalid range start IP")
	}

	endIPRaw := net.ParseIP(endIP).To4()
	if len(endIPRaw) == 0 || len(endIPRaw) > net.IPv4len {
		return IP4Range{}, errors.New("invalid range end IP")
	}

	startIPInt := binary.BigEndian.Uint32(startIPRaw)
	endIPInt := binary.BigEndian.Uint32(endIPRaw)

	// ensure valid range.
	if startIPInt > endIPInt {
		return IP4Range{}, fmt.Errorf("invalid range: %v > %v", startIP, endIP)
	}

	return IP4Range{
		startIP: IP(startIPInt),
		endIP:   IP(endIPInt),
	}, nil
}

// NewIP6RangeInt checks if the starting and ending IP range is valid or not.
//
// note: starting and ending IPs must be valid IPv6 string formats.
func NewIP6RangeInt(
	startIP string,
	endIP string,
) (IP6RangeInt, error) {
	startIPByte := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	endIPByte := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	startIPRaw := net.ParseIP(startIP).To16()
	if len(startIPRaw) == 0 {
		return IP6RangeInt{}, errors.New("invalid range start IP")
	}

	endIPRaw := net.ParseIP(endIP)
	if len(endIPRaw) == 0 {
		return IP6RangeInt{}, errors.New("invalid range end IP")
	}
	copy(startIPByte[:], []byte(startIPRaw.To16()))
	copy(endIPByte[:], []byte(endIPRaw.To16()))

	startIPInt := new(big.Int)
	endIPInt := new(big.Int)
	startIPInt.SetBytes(startIPByte[:])
	endIPInt.SetBytes(endIPByte[:])

	// ensure valid range
	if startIPInt.Cmp(endIPInt) > 0 {
		return IP6RangeInt{}, fmt.Errorf("invalid range: %v > %v", startIP, endIP)
	}
	return IP6RangeInt{
		startIP: startIPInt,
		endIP:   endIPInt,
	}, nil
}

// RandIP4Range returns a list of randomly generated IPv4 addresses within
// the range specified by `startIP` and `endIP`.
//
// note: `EvalIP4` must be called before this function as this function assumes
// `startIP` and `endIP` is a correct range.
func RandIP4Range(iprange IP4Range) (net.IP, error) {
	tmp := iprange.endIP - iprange.startIP
	if tmp == 0 {
		tmpIP := [4]byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(tmpIP[:], uint32(iprange.startIP))
		return tmpIP[:], nil
	}

	// get random IP and adjust it to fit range.
	randIP := binary.BigEndian.Uint32(RandIP4().To4())
	randIP %= (uint32(iprange.endIP) - uint32(iprange.startIP))
	randIP += uint32(iprange.startIP)
	randIPbyte := [4]byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(randIPbyte[:], randIP)
	return net.IP(randIPbyte[:]), nil
}

// RandIP4ListWrite prints a list of new randomly generated IPv4 addresses.
func RandIP4ListWrite(n int) {
	for i := 0; i < n; i++ {
		fmt.Println(RandIP4())
	}
}

// RandIP4ListWrite prints a list of new randomly generated IPv4 addresses
// within starting and IPs ending range.
func RandIP4RangeListWrite(startIP, endIP string, n int) error {
	ipRange, err := NewIP4Range(startIP, endIP)
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		ip, err := RandIP4Range(ipRange)
		if err != nil {
			return err
		}
		fmt.Println(ip)
	}
	return nil
}

// IPFromStdIP returns a new IPv4 address representation from the standard
// library's IP representation.
func IPFromStdIP(ip net.IP) IP {
	return IP(binary.BigEndian.Uint32(ip.To4()))
}

// RandIP6 returns a new randomly generated IPv6 address.
func RandIP6() net.IP {
	ip := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint64(ip[0:], rand.Uint64())
	binary.BigEndian.PutUint64(ip[8:], rand.Uint64())
	return net.IP(ip[:])
}

// RandIP6Range returns a list of randomly generated IPv6 addresses within
// the range specified by `startIP` and `endIP`.
//
// note: `EvalIP6` must be called before this function as this function assumes
// `startIP` and `endIP` is a correct range.
func RandIP6Range(ipRange IP6RangeInt) (net.IP, error) {
	randIP := new(big.Int)
	randIP.SetBytes(RandIP6())
	tmp := new(big.Int)
	tmp.Sub(ipRange.endIP, ipRange.startIP)
	if tmp.Cmp(big.NewInt(0)) <= 0 {
		randIPBytes := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		randIPBigIntBytes := ipRange.startIP.Bytes()
		copy(randIPBytes[16-len(randIPBigIntBytes):], randIPBigIntBytes)
		return net.IP(randIPBytes[:]), nil
	}
	randIP.Mod(randIP, tmp)
	randIP.Add(randIP, ipRange.startIP)
	randIPBytes := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	randIPBigIntBytes := randIP.Bytes()
	copy(randIPBytes[16-len(randIPBigIntBytes):], randIPBigIntBytes)
	return net.IP(randIPBytes[:]).To16(), nil
}

// RandIP6List returns a list of new randomly generated IPv6 addresses.
func RandIP6List(n int) []net.IP {
	ips := make([]net.IP, n)
	for i := 0; i < n; i++ {
		ips[i] = RandIP6()
	}
	return ips
}

// RandIP6ListWrite prints a list of randomly generated IPv6 addresses.
func RandIP6ListWrite(n int) {
	for i := 0; i < n; i++ {
		fmt.Println(RandIP6())
	}
}

// RandIP6ListWrite prints a list of new randomly generated IPv6 addresses
// withing starting and ending IPs range.
func RandIP6RangeListWrite(startIP, endIP string, n int) error {
	ipRange, err := NewIP6RangeInt(startIP, endIP)
	if err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		ip, err := RandIP6Range(ipRange)
		if err != nil {
			return err
		}
		fmt.Println(ip)
	}
	return nil
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
