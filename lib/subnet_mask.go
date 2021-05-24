package lib

import (
	"net"
	"strconv"
)

type SubnetMask struct {
	NetBitCnt uint32
	NetMask uint32
	HostBitCnt uint32
	HostMask uint32
	LoIP uint32
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
		mask += 1<<(32-(i+1))
	}

	return mask, ^mask
}

// LargestSubnetMaskWithinRange returns the largest subnet mask that does not
// exceed the range `[start,end]` where `start <= end`.
func LargestSubnetMaskWithinRange(
	start uint32,
	end uint32,
) SubnetMask {
	// start with a mask which is equal to CIDR `<start>/32`.
	subnetMask := SubnetMask{
		NetBitCnt: 32,
		NetMask: 0xffffffff,
		HostBitCnt: 0,
		HostMask: 0,
		LoIP: start,
		HiIP: start,
	}

	for i := 31; i >= 0; i-- {
		netMask, hostMask := NetAndHostMasks(uint32(i))
		tmpSubnetMask := SubnetMask{
			NetBitCnt: uint32(i),
			NetMask: netMask,
			HostBitCnt: 32-uint32(i),
			HostMask: hostMask,
			LoIP: start&netMask,
			HiIP: (start&netMask)|hostMask,
		}

		if tmpSubnetMask.LoIP == start && tmpSubnetMask.HiIP <= end {
			subnetMask = tmpSubnetMask
		} else {
			break
		}
	}

	return subnetMask
}

// SubnetMasksWithinRange returns a list of subnet masks which cover the full
// range `[start,end]` where `start <= end`.
func SubnetMasksWithinRange(
	start uint32,
	end uint32,
) []SubnetMask {
	// use u64 versions so we don't have overflow when doing some arithmetic.
	_start := uint64(start)
	_end := uint64(end)
	subnetMasks := make([]SubnetMask, 0)

	for _start <= _end {
		subnetMask := LargestSubnetMaskWithinRange(
			uint32(_start), uint32(_end),
		)
		subnetMasks = append(subnetMasks, subnetMask)
		_start = uint64(subnetMask.HiIP) + uint64(1)
	}

	return subnetMasks
}

// ToCIDR converts a SubnetMask to CIDR notation.
func (s SubnetMask) ToCIDR() string {
	loIPStr := net.IPv4(
		byte(s.LoIP>>24),
		byte(s.LoIP>>16),
		byte(s.LoIP>>8),
		byte(s.LoIP),
	).String()
	netBitCntStr := strconv.Itoa(int(s.NetBitCnt))
	return loIPStr+"/"+netBitCntStr
}
