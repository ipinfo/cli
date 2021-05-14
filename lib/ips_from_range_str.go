package lib

import (
	"net"
)

// IPsFromRangeStr returns a list of IPs given a range string.
// The string must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func IPsFromRangeStr(r string) ([]net.IP, error) {
	rStart, rEnd, err := IPRangeStrPartsFromRangeStr(r)
	if err != nil {
		return nil, err
	}

	return IPsFromRange(rStart, rEnd)
}
