package lib

import (
	"net"
	"strings"
)

// IPRangeStrPartsFromRangeStr returns the two IP parts (start and end) of an
// IP range string.
// The string must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func IPRangeStrPartsFromRangeStr(r string) (string, string, error) {
	idx := strings.IndexAny(r, "-,")
	if idx == -1 || idx == len(r)-1 {
		return "", "", ErrNotIPRange
	}

	rStart, rEnd := r[:idx], r[idx+1:]
	if net.ParseIP(rStart) == nil || net.ParseIP(rEnd) == nil {
		return "", "", ErrNotIPRange
	}

	return rStart, rEnd, nil
}