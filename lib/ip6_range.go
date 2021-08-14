package lib

import (
	"encoding/binary"
	"net"
)

// IP6Range represents a range of IPv6 addresses [Start, End].
type IP6Range struct {
	// Start is the first IP in the IPv6 range.
	Start IP6

	// End is the last IP in the IPv6 range.
	End IP6
}

// IP6Range returns a new IPv6 address range given a start and end IPv6
// address.
func NewIP6Range(start IP6, end IP6) IP6Range {
	return IP6Range{Start: start, End: end}
}

// IP6RangeFromCIDR returns an IP6Range given a CIDR in string form.
func IP6RangeFromCIDR(cidrStr string) (IP6Range, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return IP6Range{}, err
	}

	maskhi := binary.BigEndian.Uint64(ipnet.Mask[:8])
	masklo := binary.BigEndian.Uint64(ipnet.Mask[8:])
	starthi := binary.BigEndian.Uint64(ipnet.IP[:8])
	startlo := binary.BigEndian.Uint64(ipnet.IP[8:])
	endhi := (starthi & maskhi) | (maskhi ^ 0xffffffffffffffff)
	endlo := (startlo & masklo) | (masklo ^ 0xffffffffffffffff)

	start := NewIP6(starthi, startlo)
	end := NewIP6(endhi, endlo)
	return NewIP6Range(start, end), nil
}
