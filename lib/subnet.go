package lib

import (
	"net"
	"strconv"
)

// Subnet is the representation of a subnet.
type Subnet struct {
	// NetBitCnt is the number of bits in the network part of the subnet.
	NetBitCnt uint32

	// NetMask is the subnet mask of the network part of the subnet.
	NetMask uint32

	// HostBitCnt is the number of bits in the host part of the subnet.
	HostBitCnt uint32

	// HostMask is the subnet mask of the host part of the subnet.
	HostMask uint32

	// LoIP is the big-endian representation of the lowest IP in the subnet.
	LoIP uint32

	// HiIP is the big-endian representation of the highest IP in the subnet.
	HiIP uint32
}

// NetAndHostMasks returns a network and host masks where the `size`
// most-significant bits are set to 1 and the rest set to 0 in the network
// mask, and the host mask is the bitwise-negation of the network mask.
func NetAndHostMasks(size uint32) (uint32, uint32) {
	if size > 32 {
		size = 32
	}

	var mask uint32 = 0
	for i := uint32(0); i < size; i++ {
		mask += 1 << (32 - (i + 1))
	}

	return mask, ^mask
}

// LargestSubnetWithinRange returns the largest subnet mask that does not
// exceed the range `[start,end]` where `start <= end`.
func LargestSubnetWithinRange(
	start uint32,
	end uint32,
) Subnet {
	// start with a subnet which is equal to CIDR `<start>/32`.
	subnet := Subnet{
		NetBitCnt:  32,
		NetMask:    0xffffffff,
		HostBitCnt: 0,
		HostMask:   0,
		LoIP:       start,
		HiIP:       start,
	}

	for i := 31; i >= 0; i-- {
		netMask, hostMask := NetAndHostMasks(uint32(i))
		tmpSubnet := Subnet{
			NetBitCnt:  uint32(i),
			NetMask:    netMask,
			HostBitCnt: 32 - uint32(i),
			HostMask:   hostMask,
			LoIP:       start & netMask,
			HiIP:       (start & netMask) | hostMask,
		}

		if tmpSubnet.LoIP == start && tmpSubnet.HiIP <= end {
			subnet = tmpSubnet
		} else {
			break
		}
	}

	return subnet
}

// SubnetsWithinRange returns a list of subnet masks which cover the full
// range `[start,end]` where `start <= end`.
func SubnetsWithinRange(
	start uint32,
	end uint32,
) []Subnet {
	// use u64 versions so we don't have overflow when doing some arithmetic.
	_start := uint64(start)
	_end := uint64(end)
	subnets := make([]Subnet, 0)

	for _start <= _end {
		subnet := LargestSubnetWithinRange(uint32(_start), uint32(_end))
		subnets = append(subnets, subnet)
		_start = uint64(subnet.HiIP) + uint64(1)
	}

	return subnets
}

// ToCIDR converts a Subnet to CIDR notation.
func (s Subnet) ToCIDR() string {
	loIPStr := net.IPv4(
		byte(s.LoIP>>24),
		byte(s.LoIP>>16),
		byte(s.LoIP>>8),
		byte(s.LoIP),
	).String()
	netBitCntStr := strconv.Itoa(int(s.NetBitCnt))
	return loIPStr + "/" + netBitCntStr
}
