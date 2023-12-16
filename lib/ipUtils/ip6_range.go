package ipUtils

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

// LargestIP6Subnet returns the largest subnet mask within this range.
func (r IP6Range) LargestIP6Subnet() IP6Subnet {
	// start with a subnet which is equal to CIDR `<start>/128`.
	start := r.Start
	end := r.End
	subnet := IP6Subnet{
		NetBitCnt:  128,
		NetMask:    MaxU128,
		HostBitCnt: 0,
		HostMask:   ZeroU128,
		LoIP:       start,
		HiIP:       start,
	}

	for i := 127; i >= 0; i-- {
		netMask, hostMask := NetAndHostMasks6(uint32(i))
		tmpSubnet := IP6Subnet{
			NetBitCnt:  uint32(i),
			NetMask:    netMask,
			HostBitCnt: 128 - uint32(i),
			HostMask:   hostMask,
			LoIP:       IP6FromU128(start.N.And(netMask)),
			HiIP:       IP6FromU128(start.N.And(netMask).Or(hostMask)),
		}

		if tmpSubnet.LoIP.Eq(start) && tmpSubnet.HiIP.Lte(end) {
			subnet = tmpSubnet
		} else {
			break
		}
	}

	return subnet
}

// ToIP6Subnets returns a list of subnet masks which cover the full IP range.
func (r IP6Range) ToIP6Subnets() []IP6Subnet {
	start := r.Start.N
	end := r.End.N
	subnets := make([]IP6Subnet, 0)

	for start.Lte(end) {
		subnet := NewIP6Range(
			IP6FromU128(start),
			IP6FromU128(end),
		).LargestIP6Subnet()
		subnets = append(subnets, subnet)
		start = subnet.HiIP.N.AddOne()
		if start.IsZero() {
			break
		}
	}

	return subnets
}

// ToCIDRs returns a list of CIDR strings which cover the full IP range.
func (r IP6Range) ToCIDRs() []string {
	subnets := r.ToIP6Subnets()
	cidrStrs := make([]string, len(subnets))
	for i, subnet := range subnets {
		cidrStrs[i] = subnet.ToCIDR()
	}
	return cidrStrs
}
