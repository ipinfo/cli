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

// CIDRToIPSubnet converts a CIDR notation to IPSubnet.
func CIDRToIPSubnet(cidr string) (IPSubnet, error) {
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

// SubnetBitShift returns a list of IPSubnet after shifting number of `bits`.
func (s IPSubnet) SubnetBitShift(bits int) ([]IPSubnet, error) {
	ipsubnets := make([]IPSubnet, 1<<bits)
	_, network, err := net.ParseCIDR(s.ToCIDR())
	if err != nil {
		return nil, err
	}
	subnets, err := NetworkBitShift(network, bits)
	if err != nil {
		return nil, err
	}
	for i, subnet := range subnets {
		ipsubnet, err := CIDRToIPSubnet(fmt.Sprint(subnet))
		if err != nil {
			return nil, err
		}
		ipsubnets[i] = ipsubnet
	}

	return ipsubnets, nil
}

// NetworkBitShift returns a list of *net.IPNet after shifting the bits number
// of on `bits` on network `*net.IPNet`.
func NetworkBitShift(network *net.IPNet, bits int) ([]*net.IPNet, error) {
	subnets := make([]*net.IPNet, 1<<bits)
	ones, _ := network.Mask.Size()
	hostBits := (32 - ones) - bits
	newMask := net.CIDRMask(int(ones+bits), 32)
	for i := range subnets {
		ip := binary.BigEndian.Uint32(network.IP) + uint32(i*(1<<uint(hostBits)))
		ip4 := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip4, ip)
		subnets[i] = &net.IPNet{
			IP:   ip4,
			Mask: newMask,
		}
	}

	return subnets, nil
}
