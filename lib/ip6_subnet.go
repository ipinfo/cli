package lib

import (
	"strconv"
)

// IP6Subnet is the representation of a IPv6 subnet.
type IP6Subnet struct {
	// NetBitCnt is the number of bits in the network part of the subnet.
	NetBitCnt uint32

	// NetMask is the subnet mask of the network part of the subnet.
	NetMask U128

	// HostBitCnt is the number of bits in the host part of the subnet.
	HostBitCnt uint32

	// HostMask is the subnet mask of the host part of the subnet.
	HostMask U128

	// LoIP is the big-endian representation of the lowest IP in the subnet.
	LoIP IP6

	// HiIP is the big-endian representation of the highest IP in the subnet.
	HiIP IP6
}

// IP6CIDR is the representation of a IPv6 subnet in CIDR form.
type IP6CIDR string

// NetAndHostMasks returns network and host masks where the `size`
// most-significant bits are set to 1 and the rest set to 0 in the network
// mask, and the host mask is the bitwise-negation of the network mask.
func NetAndHostMasks6(size uint32) (U128, U128) {
	if size > 128 {
		size = 128
	}

	mask := U128BitMasks[size]
	return mask, mask.Not()
}

// ToCIDR converts a IPSubnet to CIDR notation.
func (s IP6Subnet) ToCIDR() string {
	netBitCntStr := strconv.Itoa(int(s.NetBitCnt))
	return s.LoIP.String() + "/" + netBitCntStr
}

// CIDRsFromIP6RangeStrRaw returns a list of CIDR strings which cover the full
// range specified in the IP range string `rStr`.
//
// `rStr` must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func CIDRsFromIP6RangeStrRaw(rStr string) ([]string, error) {
	r, err := IP6RangeStrFromStr(rStr)
	if err != nil {
		return nil, err
	}

	return r.ToCIDRs(), nil
}
