package lib

import (
	"encoding/binary"
	"net"
)

// CIDRsFromIPRange returns a list of CIDR strings which cover the full range
// `[start,end]` where `start <= end`.
func CIDRsFromIPRange(
	start uint32,
	end uint32,
) []string {
	subnets := SubnetsWithinRange(start, end)
	cidrStrs := make([]string, len(subnets))
	for i, subnet := range subnets {
		cidrStrs[i] = subnet.ToCIDR()
	}
	return cidrStrs
}

// CIDRsFromIPRangeStr returns a list of CIDR strings which cover the full
// range specified in the IP range string `r`.
// The string must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func CIDRsFromIPRangeStr(
	r string,
) ([]string, error) {
	startStr, endStr, err := IPRangeStrFromRangeStr(r)
	if err != nil {
		return nil, err
	}

	start := binary.BigEndian.Uint32(net.ParseIP(startStr).To4())
	end := binary.BigEndian.Uint32(net.ParseIP(endStr).To4())
	if start <= end {
		return CIDRsFromIPRange(start, end), nil
	}

	cidrStrs := CIDRsFromIPRange(end, start)
	ReverseSliceString(cidrStrs)
	return cidrStrs, nil
}
