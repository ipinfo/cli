package lib

import (
	"encoding/binary"
	"math/big"
	"net"
)

func ipAdd(input string, delta int) net.IP {
	ip := net.ParseIP(input)
	if ip.To4() != nil {
		ipInt := ipToUint32(ip)
		newIPInt := ipInt + uint32(delta)
		newIP := uint32ToIP(newIPInt)
		return newIP
	} else {
		ipInt := ipToBigInt(ip)
		deltaBigInt := new(big.Int).SetInt64(int64(delta))
		newIPInt := new(big.Int).Add(ipInt, deltaBigInt)
		adjustedIPInt := adjustIPBigInt(newIPInt)
		newIP := bigIntToIP(adjustedIPInt)
		return newIP
	}
}

func ipToUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

func uint32ToIP(ipInt uint32) net.IP {
	ip := make(net.IP, net.IPv4len)
	binary.BigEndian.PutUint32(ip, ipInt)
	return ip
}

func ipToBigInt(ip net.IP) *big.Int {
	ipInt := new(big.Int)
	ipInt.SetBytes(ip)
	return ipInt
}

func bigIntToIP(ipInt *big.Int) net.IP {
	ipIntBytes := ipInt.Bytes()
	ip := make(net.IP, net.IPv6len)
	copyLength := len(ipIntBytes)
	if copyLength > net.IPv6len {
		copyLength = net.IPv6len
	}
	copy(ip[net.IPv6len-copyLength:], ipIntBytes)
	return ip
}

func adjustIPBigInt(ipInt *big.Int) *big.Int {
	if ipInt.Cmp(maxIPv6BigInt) > 0 {
		ipInt.Sub(ipInt, maxIPv6BigInt)
		ipInt.Sub(ipInt, big.NewInt(1))
	} else if ipInt.Cmp(big.NewInt(0)) < 0 {
		ipInt.Add(ipInt, maxIPv6BigInt)
		ipInt.Add(ipInt, big.NewInt(1))
	}
	return ipInt
}

var (
	maxIPv6BigInt, _ = new(big.Int).SetString("340282366920938463463374607431768211455", 10) // 2^128 - 1
)
