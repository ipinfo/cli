package ipUtils

import (
	"net"
	"strings"
)

// IP6RangeStr represents a range of IPv6 addresses [Start, End] in string
// form.
type IP6RangeStr struct {
	// Start is the first IP in the IPv6 range.
	Start string

	// End is the last IP in the IPv6 range.
	End string
}

// NewIP6RangeStr returns a new IP range string given a start & end IP.
func NewIP6RangeStr(start string, end string) IP6RangeStr {
	return IP6RangeStr{Start: start, End: end}
}

// StrIsIP6RangeStr checks whether a string is an IPv6 range string.
//
// The string must be of any of these forms to be considered an IP range:
//
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func StrIsIP6RangeStr(r string) bool {
	_, err := IP6RangeStrFromStr(r)
	return err == nil
}

// IP6RangeStrFromStr returns the two IPv6 parts (start and end) of an IPv6
// range string.
//
// The string must be of any of these forms:
//
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func IP6RangeStrFromStr(r string) (IP6RangeStr, error) {
	idx := strings.IndexAny(r, "-,")
	if idx == -1 || idx == len(r)-1 {
		return IP6RangeStr{}, ErrNotIP6Range
	}

	rStart, rEnd := r[:idx], r[idx+1:]
	if net.ParseIP(rStart) == nil || net.ParseIP(rEnd) == nil {
		return IP6RangeStr{}, ErrNotIP6Range
	}

	return NewIP6RangeStr(rStart, rEnd), nil
}

// IP6RangeStrFromCIDR returns the start and end IPv6 strings of a CIDR.
func IP6RangeStrFromCIDR(cidrStr string) (IP6RangeStr, error) {
	r, err := IP6RangeFromCIDR(cidrStr)
	if err != nil {
		return IP6RangeStr{}, err
	}

	return NewIP6RangeStr(r.Start.String(), r.End.String()), nil
}

// ToIP6Range converts the string form of the range into a numerical form.
func (r IP6RangeStr) ToIP6Range() IP6Range {
	start, _ := IP6FromStdIP(net.ParseIP(r.Start).To16())
	end, _ := IP6FromStdIP(net.ParseIP(r.End).To16())
	return NewIP6Range(start, end)
}

// ToCIDRs returns a list of CIDR strings which cover the full range specified
// in the IP range string `r`.
func (rStr IP6RangeStr) ToCIDRs() []string {
	rev := false
	r := rStr.ToIP6Range()
	if r.Start.Gt(r.End) {
		rev = true
		tmp := r.Start
		r.Start = r.End
		r.End = tmp
	}
	cidrStrs := r.ToCIDRs()
	if rev {
		StringSliceRev(cidrStrs)
	}
	return cidrStrs
}

// String returns the IPv6 range string as `<start>-<end>`.
func (r IP6RangeStr) String() string {
	return r.Start + "-" + r.End
}

// StringDelim is the same as String but allows a custom delimiter.
func (r IP6RangeStr) StringDelim(d string) string {
	return r.Start + d + r.End
}
