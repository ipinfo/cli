package ipUtils

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"net"
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
//
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func CIDRsFromIP6RangeStrRaw(rStr string) ([]string, error) {
	r, err := IP6RangeStrFromStr(rStr)
	if err != nil {
		return nil, err
	}

	return r.ToCIDRs(), nil
}

// IP6SubnetFromCidr converts a CIDR notation to IP6Subnet.
func IP6SubnetFromCidr(cidr string) (IP6Subnet, error) {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return IP6Subnet{}, err
	}

	ones, _ := network.Mask.Size()
	netMask, hostMask := NetAndHostMasks6(uint32(ones))
	starthi := binary.BigEndian.Uint64(network.IP[:8])
	startlo := binary.BigEndian.Uint64(network.IP[8:])
	start := NewIP6(starthi, startlo)
	ip6subnet := IP6Subnet{
		HostBitCnt: uint32(128 - ones),
		HostMask:   hostMask,
		NetBitCnt:  uint32(ones),
		NetMask:    netMask,
		LoIP:       IP6FromU128(start.N.And(netMask)),
		HiIP:       IP6FromU128(start.N.And(netMask).Or(hostMask)),
	}

	return ip6subnet, nil
}

// SplitCIDR returns a list of smaller IPSubnet after splitting a larger CIDR
// into `split`.
func (s IP6Subnet) SplitCIDR(split int) ([]IP6Subnet, error) {
	bitshifts := int(uint32(split) - s.NetBitCnt)
	if bitshifts < 0 || bitshifts > 128 || int(s.NetBitCnt)+bitshifts > 128 {
		return nil, fmt.Errorf("wrong split string")
	}

	hostBits := (128 - s.NetBitCnt) - uint32(bitshifts)
	netMask, hostMask := NetAndHostMasks6(uint32(split))
	subnetCount := math.Pow(2, float64(bitshifts))
	subnetCountBig := big.NewFloat(subnetCount)
	hostCount := math.Pow(2, float64(hostBits))
	hostCountbig := big.NewFloat(hostCount)

	var ipsubnets []IP6Subnet
	for i := big.NewFloat(0); i.Cmp(subnetCountBig) < 0; i.Add(i, big.NewFloat(1)) {
		// calculating new LoIP by `LoIP + i*(hostCount)`
		hostCountMul := new(big.Float)
		hostCountMul.Mul(i, hostCountbig)
		newIP := new(big.Int)
		newIP.SetBytes(s.LoIP.To16ByteSlice())
		hostCountAdd := new(big.Int)
		result := new(big.Int)
		hostCountMul.Int(result)
		hostCountAdd.Add(newIP, result)

		// converting `bigint` to `IP6`
		ipArr := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		copy(ipArr[16-len(hostCountAdd.Bytes()):], hostCountAdd.Bytes())
		startIP := NewIP6(binary.BigEndian.Uint64(ipArr[:8]), binary.BigEndian.Uint64(ipArr[8:]))

		subnet := IP6Subnet{
			HostBitCnt: uint32(128 - split),
			HostMask:   hostMask,
			NetBitCnt:  uint32(split),
			LoIP:       IP6FromU128(startIP.N.And(netMask)),
			HiIP:       IP6FromU128(startIP.N.And(netMask).Or(hostMask)),
		}
		ipsubnets = append(ipsubnets, subnet)
	}

	return ipsubnets, nil
}
