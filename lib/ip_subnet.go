package lib

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

// IPSubnet is the representation of a IPv4 subnet.
type IPSubnet struct {
	// NetBitCnt is the number of bits in the network part of the subnet.
	NetBitCnt uint32

	// NetMask is the subnet mask of the network part of the subnet.
	NetMask uint32

	// HostBitCnt is the number of bits in the host part of the subnet.
	HostBitCnt uint32

	// HostMask is the subnet mask of the host part of the subnet.
	HostMask uint32

	// LoIP is the big-endian representation of the lowest IP in the subnet.
	LoIP IP

	// HiIP is the big-endian representation of the highest IP in the subnet.
	HiIP IP
}

// IPCIDR is the representation of a IPv4 subnet in CIDR form.
type IPCIDR string

// NetAndHostMasks returns network and host masks where the `size`
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

// ToCIDR converts a IPSubnet to CIDR notation.
func (s IPSubnet) ToCIDR() string {
	netBitCntStr := strconv.Itoa(int(s.NetBitCnt))
	return s.LoIP.String() + "/" + netBitCntStr
}

// CIDRsFromIPRangeStrRaw returns a list of CIDR strings which cover the full
// range specified in the IP range string `rStr`.
//
// `rStr` must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func CIDRsFromIPRangeStrRaw(rStr string) ([]string, error) {
	r, err := IPRangeStrFromStr(rStr)
	if err != nil {
		return nil, err
	}

	return r.ToCIDRs(), nil
}

// IPSubnetFromCidr converts a CIDR notation to IPSubnet.
func IPSubnetFromCidr(cidr string) (IPSubnet, error) {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return IPSubnet{}, err
	}

	ones, _ := network.Mask.Size()
	netMask, hostMask := NetAndHostMasks(uint32(ones))
	start := binary.BigEndian.Uint32(network.IP)
	ipsubnet := IPSubnet{
		HostBitCnt: uint32(32 - ones),
		HostMask:   hostMask,
		NetBitCnt:  uint32(ones),
		NetMask:    netMask,
		LoIP:       IP(uint32(start) & netMask),
		HiIP:       IP((uint32(start) & netMask) | hostMask),
	}

	return ipsubnet, nil
}

// SplitCIDR returns a list of smaller IPSubnet after splitting a larger CIDR
// into `split`.
func (s IPSubnet) SplitCIDR(split int) ([]IPSubnet, error) {
	bitshifts := int(uint32(split) - s.NetBitCnt)
	if bitshifts < 0 || bitshifts > 31 || int(s.NetBitCnt)+bitshifts > 32 {
		return nil, fmt.Errorf("wrong split string")
	}

	hostBits := (32 - s.NetBitCnt) - uint32(bitshifts)
	netMask, hostMask := NetAndHostMasks(uint32(split))
	ipsubnets := make([]IPSubnet, 1<<bitshifts)
	for i := range ipsubnets {
		start := uint32(s.LoIP) + uint32(i*(1<<hostBits))
		ipsubnets[i] = IPSubnet{
			HostBitCnt: uint32(32 - split),
			HostMask:   hostMask,
			NetBitCnt:  uint32(split),
			LoIP:       IP(uint32(start) & netMask),
			HiIP:       IP((uint32(start) & netMask) | hostMask),
		}
	}

	return ipsubnets, nil
}
