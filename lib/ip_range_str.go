package lib

import (
	"encoding/binary"
	"net"
	"strings"
)

// IPRangeStr represents a range of IPv4 addresses [Start, End] in string form.
type IPRangeStr struct {
	// Start is the first IP in the IPv4 range.
	Start string

	// End is the last IP in the IPv4 range.
	End string
}

// NewIPRangeStr returns a new IP range string given a start & end IP.
func NewIPRangeStr(start string, end string) IPRangeStr {
	return IPRangeStr{Start: start, End: end}
}

// StrIsIPRangeStr checks whether a string is an IP range.
//
// The string must be of any of these forms to be considered an IP range:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func StrIsIPRangeStr(r string) bool {
	_, err := IPRangeStrFromStr(r)
	return err == nil
}

// IPRangeStrFromStr returns the two IP parts (start and end) of an IP
// range string.
//
// The string must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func IPRangeStrFromStr(r string) (IPRangeStr, error) {
	idx := strings.IndexAny(r, "-,")
	if idx == -1 || idx == len(r)-1 {
		return IPRangeStr{}, ErrNotIPRange
	}

	rStart, rEnd := r[:idx], r[idx+1:]
	if net.ParseIP(rStart) == nil || net.ParseIP(rEnd) == nil {
		return IPRangeStr{}, ErrNotIPRange
	}

	return NewIPRangeStr(rStart, rEnd), nil
}

// IPRangeStrFromCIDR returns the start and end IP strings of a CIDR.
func IPRangeStrFromCIDR(cidrStr string) (IPRangeStr, error) {
	r, err := IPRangeFromCIDR(cidrStr)
	if err != nil {
		return IPRangeStr{}, err
	}

	return NewIPRangeStr(r.Start.String(), r.End.String()), nil
}

// ToIPRange converts the string form of the range into a numerical form.
func (r IPRangeStr) ToIPRange() IPRange {
	start := binary.BigEndian.Uint32(net.ParseIP(r.Start).To4())
	end := binary.BigEndian.Uint32(net.ParseIP(r.End).To4())
	if start <= end {
		return NewIPRange(IP(start), IP(end))
	}
	return NewIPRange(IP(end), IP(start))
}

// ToCIDRs returns a list of CIDR strings which cover the full range specified
// in the IP range string `r`.
func (r IPRangeStr) ToCIDRs() []string {
	cidrStrs := r.ToIPRange().ToCIDRs()
	if r.Start > r.End {
		StringSliceRev(cidrStrs)
	}

	return cidrStrs
}

// String returns the IP range string as `<start>-<end>`.
func (r IPRangeStr) String() string {
	return r.Start+"-"+r.End
}

// StringDelim is the same as String but allows a custom delimiter.
func (r IPRangeStr) StringDelim(d string) string {
	return r.Start+d+r.End
}
