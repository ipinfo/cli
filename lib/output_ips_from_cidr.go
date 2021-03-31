package lib

import (
	"encoding/binary"
	"fmt"
	"net"
)

// OutputIPsFromCIDR is the same as IPsFromCIDR with O(1) memory by discarding
// IPs after printing.
func OutputIPsFromCIDR(cidrStr string) error {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return err
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)
	end := (start & mask) | (mask ^ 0xffffffff)

	for i := start; i <= end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		fmt.Println(ip)
	}

	return nil
}
