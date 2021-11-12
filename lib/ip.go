package lib

import (
	"encoding/binary"
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

// RandIP4Range returns a list of randomly generated IPv4 addresses within
// starting and ending IPrange.

// note: starting and ending IPs should be in format of "1.1.1.1"
func RandIP4Range(startIP, EndIP string) (net.IP, error) {
	StartIPRaw := net.ParseIP(startIP).To4()
	if len(StartIPRaw) == 0 || len(StartIPRaw) > net.IPv4len {
		return nil, fmt.Errorf("range is Invalid")
	}
	StartIPInt := binary.BigEndian.Uint32(StartIPRaw.To4())
	EndIPRaw := net.ParseIP(EndIP).To4()
	if len(EndIPRaw) == 0 || len(EndIPRaw) > net.IPv4len {
		return nil, fmt.Errorf("range is Invalid")
	}
	EndIPInt := binary.BigEndian.Uint32(EndIPRaw.To4())
	if StartIPInt > EndIPInt {
		return nil, fmt.Errorf("range is Invalid")
	}
	rangeIP := binary.BigEndian.Uint32(RandIP4().To4())
	temp := EndIPInt - StartIPInt
	if temp <= 0 {
		return nil, fmt.Errorf("range is Invalid")
	}
	rangeIP %= (EndIPInt - StartIPInt)
	rangeIP += StartIPInt
	rangeIPbyte := [4]byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(rangeIPbyte[:], rangeIP)
	return net.IP(rangeIPbyte[:]), nil
}

// RandIP4ListWrite prints a list of new randomly generated IPv4 addresses.
func RandIP4ListWrite(n int) {
	for i := 0; i < n; i++ {
		fmt.Println(RandIP4())
	}
}

// RandIP4ListWrite prints a list of new randomly generated IPv4 addresses
// within starting and IPs ending range.
func RandIP4RangeListWrite(min, max string, n int) error {
	for i := 0; i < n; i++ {
		ip, err := RandIP4Range(min, max)
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
// starting and ending IPrange.
func RandIP6Range(StartIP, EndIP string) (net.IP, error) {
	StartIPByte := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	EndIPByte := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	StartIPRaw := net.ParseIP(StartIP)
	EndIPRaw := net.ParseIP(EndIP)
	copy(StartIPByte[:], []byte(StartIPRaw.To16()))
	copy(EndIPByte[:], []byte(EndIPRaw.To16()))

	min_ip := new(big.Int)
	max_ip := new(big.Int)
	min_ip.SetBytes(StartIPByte[:])
	max_ip.SetBytes(EndIPByte[:])
	if min_ip.Cmp(max_ip) > 0 {
		return nil, fmt.Errorf("range is Invalid")
	}
	randIP := new(big.Int)
	randIP.SetBytes(RandIP6())
	tmp := new(big.Int)
	tmp.Sub(max_ip, min_ip)
	if tmp.Cmp(big.NewInt(0)) <= 0 {
		return nil, fmt.Errorf("range is Invalid")
	}
	randIP.Mod(randIP, tmp)
	randIP.Add(randIP, min_ip)
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
func RandIP6RangeListWrite(min, max string, n int) error {
	for i := 0; i < n; i++ {
		ip, err := RandIP6Range(min, max)
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
