package iputil

import (
	"encoding/binary"
	"net"
)

// IPRange represents a range of IPv4 addresses [Start, End].
type IPRange struct {
	// Start is the first IP in the IPv4 range.
	Start IP

	// End is the last IP in the IPv4 range.
	End IP
}

// NewIPRange returns a new IP range given a start and end IP.
func NewIPRange(start IP, end IP) IPRange {
	return IPRange{Start: start, End: end}
}

// IPRangeFromCIDR returns the start and end IPs in big endian byte order of a
// CIDR in string form.
func IPRangeFromCIDR(cidrStr string) (IPRange, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return IPRange{}, err
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)
	end := (start & mask) | (mask ^ 0xffffffff)
	return NewIPRange(IP(start), IP(end)), nil
}

// LargestIPSubnet returns the largest subnet mask within this range.
func (r IPRange) LargestIPSubnet() IPSubnet {
	// start with a subnet which is equal to CIDR `<start>/32`.
	start := r.Start
	end := r.End
	subnet := IPSubnet{
		NetBitCnt:  32,
		NetMask:    0xffffffff,
		HostBitCnt: 0,
		HostMask:   0,
		LoIP:       start,
		HiIP:       start,
	}

	for i := 31; i >= 0; i-- {
		netMask, hostMask := NetAndHostMasks(uint32(i))
		tmpSubnet := IPSubnet{
			NetBitCnt:  uint32(i),
			NetMask:    netMask,
			HostBitCnt: 32 - uint32(i),
			HostMask:   hostMask,
			LoIP:       IP(uint32(start) & netMask),
			HiIP:       IP((uint32(start) & netMask) | hostMask),
		}

		if tmpSubnet.LoIP == start && tmpSubnet.HiIP <= end {
			subnet = tmpSubnet
		} else {
			break
		}
	}

	return subnet
}

// ToIPSubnets returns a list of subnet masks which cover the full IP range.
func (r IPRange) ToIPSubnets() []IPSubnet {
	// use u64 versions so we don't have overflow when doing some arithmetic.
	start := uint64(r.Start)
	end := uint64(r.End)
	subnets := make([]IPSubnet, 0)

	for start <= end {
		subnet := NewIPRange(IP(start), IP(end)).LargestIPSubnet()
		subnets = append(subnets, subnet)
		start = uint64(subnet.HiIP) + uint64(1)
	}

	return subnets
}

// ToCIDRs returns a list of CIDR strings which cover the full IP range.
func (r IPRange) ToCIDRs() []string {
	subnets := r.ToIPSubnets()
	cidrStrs := make([]string, len(subnets))
	for i, subnet := range subnets {
		cidrStrs[i] = subnet.ToCIDR()
	}
	return cidrStrs
}
