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
func RandIP4(noBogon bool) net.IP {
	ipBytes := [4]byte{0, 0, 0, 0}
IP:
	binary.BigEndian.PutUint32(ipBytes[:], rand.Uint32())
	ip := binary.BigEndian.Uint32(ipBytes[:])
	if noBogon && IsBogonIP4(ip) {
		goto IP
	}
	return net.IPv4(ipBytes[0], ipBytes[1], ipBytes[2], ipBytes[3])
}

// RandIP4List returns a list of new randomly generated IPv4 addresses.
func RandIP4List(n int, noBogon bool) []net.IP {
	ips := make([]net.IP, n)
	for i := 0; i < n; i++ {
		ips[i] = RandIP4(noBogon)
	}
	return ips
}

// NewIP4Range returns a new IP4Range given the input start & end IPs.
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

// NewIP6RangeInt returns a new IP6RangeInt given the input start & end IPs.
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

	endIPRaw := net.ParseIP(endIP).To16()
	if len(endIPRaw) == 0 {
		return IP6RangeInt{}, errors.New("invalid range end IP")
	}
	copy(startIPByte[:], []byte(startIPRaw))
	copy(endIPByte[:], []byte(endIPRaw))

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
// the range specified by `IP4Range`.
//
// note: `NewIP4Range` must be called before this function as this function
// assumes `IP4Range` provided to the function is a correct range.
func RandIP4Range(iprange IP4Range, noBogon bool) (net.IP, error) {
	tmp := iprange.endIP - iprange.startIP
	if tmp == 0 {
		tmpIP := [4]byte{0, 0, 0, 0}
		binary.BigEndian.PutUint32(tmpIP[:], uint32(iprange.startIP))
		return tmpIP[:], nil
	}

	// get random IP and adjust it to fit range.
	var randIP uint64
	randIP = uint64(binary.BigEndian.Uint32(RandIP4(noBogon).To4()))
	randIP %= (uint64(iprange.endIP) - uint64(iprange.startIP)) + 1
	randIP += uint64(iprange.startIP)
	randIPbyte := [4]byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(randIPbyte[:], uint32(randIP))
	return net.IP(randIPbyte[:]), nil
}

// RandIP4ListWrite prints a list of randomly generated IPv4 addresses.
func RandIP4ListWrite(n int, noBogon bool) {
	for i := 0; i < n; i++ {
		fmt.Println(RandIP4(noBogon))
	}
}

// RandIP4RangeListWrite prints a list of randomly generated IPv4 addresses.
// `startIP` and `endIP` are the start & end IPs to generate IPs between.
// `n` is the number of IPs to generate.
// `noBogon`, if true, will ensure that none of the generated IPs are bogons.
// `unique`, if true, will ensure every IP generated is unique.
func RandIP4RangeListWrite(
	startIP string,
	endIP string,
	n int,
	noBogon bool,
	unique bool,
) error {
	ipRange, err := NewIP4Range(startIP, endIP)
	if err != nil {
		return err
	}
	if unique {
		// ensure range is larger than number of IPs to generate.
		if uint32(ipRange.endIP-ipRange.startIP+1) < uint32(n) {
			return errors.New("range is too small for unique IPs")
		}

		uniqueIP := make(map[uint32]struct{})
		for i := 0; i < n; i++ {
		unique:
			ip, err := RandIP4Range(ipRange, noBogon)
			if err != nil {
				return err
			}
			// does IP already exist? if so try again.
			ipInt := binary.BigEndian.Uint32(ip)
			if _, ok := uniqueIP[ipInt]; ok {
				goto unique
			}
			uniqueIP[ipInt] = struct{}{}
			fmt.Println(ip)
		}
	} else {
		for i := 0; i < n; i++ {
			ip, err := RandIP4Range(ipRange, noBogon)
			if err != nil {
				return err
			}
			fmt.Println(ip)
		}
	}
	return nil
}

// IPFromStdIP returns a new IPv4 address representation from the standard
// library's IP representation.
func IPFromStdIP(ip net.IP) IP {
	return IP(binary.BigEndian.Uint32(ip.To4()))
}

// RandIP6 returns a new randomly generated IPv6 address.
func RandIP6(noBogon bool) net.IP {
	ip := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
IP:
	binary.BigEndian.PutUint64(ip[0:], rand.Uint64())
	binary.BigEndian.PutUint64(ip[8:], rand.Uint64())
	var ip6 U128
	ip6.Hi = binary.BigEndian.Uint64(ip[0:])
	ip6.Lo = binary.BigEndian.Uint64(ip[8:])
	if noBogon && IsBogonIP6(ip6) {
		goto IP
	}
	return net.IP(ip[:])
}

// RandIP6Range returns a list of randomly generated IPv6 addresses within
// the range specified by `IP6RangeInt`.
//
// note: `NewIP6RangeInt` must be called before this function as this function
// assumes `IP6RangeInt` provided to the function is a correct range.
func RandIP6Range(ipRange IP6RangeInt, noBogon bool) net.IP {
	randIP := new(big.Int)
	randIP.SetBytes(RandIP6(noBogon))
	tmp := new(big.Int)
	tmp.Sub(ipRange.endIP, ipRange.startIP)
	if tmp.Cmp(big.NewInt(0)) == 0 {
		// convert multi-precision byte form into 16-byte IPv6 form.
		randIPBytes := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		randIPBigIntBytes := ipRange.startIP.Bytes()
		copy(randIPBytes[16-len(randIPBigIntBytes):], randIPBigIntBytes)
		return net.IP(randIPBytes[:]).To16()
	}
	tmp.Add(tmp, big.NewInt(1))
	randIP.Mod(randIP, tmp)
	randIP.Add(randIP, ipRange.startIP)
	// convert multi-precision byte form into 16-byte IPv6 form.
	randIPBytes := [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	randIPBigIntBytes := randIP.Bytes()
	copy(randIPBytes[16-len(randIPBigIntBytes):], randIPBigIntBytes)
	return net.IP(randIPBytes[:]).To16()
}

// RandIP6List returns a list of new randomly generated IPv6 addresses.
func RandIP6List(n int, noBogon bool) []net.IP {
	ips := make([]net.IP, n)
	for i := 0; i < n; i++ {
		ips[i] = RandIP6(noBogon)
	}
	return ips
}

// RandIP6ListWrite prints a list of randomly generated IPv6 addresses.
func RandIP6ListWrite(n int, noBogon bool) {
	for i := 0; i < n; i++ {
		fmt.Println(RandIP6(noBogon))
	}
}

// RandIP6RangeListWrite prints a list of randomly generated IPv6 addresses.
// `startIP` and `endIP` are the start & end IPs to generate IPs between.
// `n` is the number of IPs to generate.
// `noBogon`, if true, will ensure that none of the generated IPs are bogons.
// `unique`, if true, will ensure every IP generated is unique.
func RandIP6RangeListWrite(
	startIP string,
	endIP string,
	n int,
	noBogon bool,
	unique bool,
) error {
	ipRange, err := NewIP6RangeInt(startIP, endIP)
	if err != nil {
		return err
	}
	if unique {
		// ensure range is larger than number of IPs to generate.
		tmp := new(big.Int)
		tmp.Sub(ipRange.endIP, ipRange.startIP)
		tmp.Add(tmp, big.NewInt(1))
		count := new(big.Int).SetUint64(uint64(n))
		if tmp.Cmp(count) < 0 {
			return errors.New("range is too small for unique IPs")
		}

		uniqueIP := make(map[IP6]struct{})
		for i := 0; i < n; i++ {
		unique:
			ip := RandIP6Range(ipRange, noBogon)
			var ipInt IP6
			ipInt.N.Hi = binary.BigEndian.Uint64(ip[0:])
			ipInt.N.Lo = binary.BigEndian.Uint64(ip[8:])
			// does IP already exist? if so try again.
			if _, ok := uniqueIP[ipInt]; ok {
				goto unique
			}
			uniqueIP[ipInt] = struct{}{}
			fmt.Println(ip)
		}
	} else {
		for i := 0; i < n; i++ {
			fmt.Println(RandIP6Range(ipRange, noBogon))
		}
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
