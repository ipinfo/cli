package lib

import (
	"net"
	"encoding/binary"
	"fmt"
)

// Same as ipsFromRange with O(1) memory by discarding IPs after printing.
func OutputIPsFromRange(ipStrStart string, ipStrEnd string) error {
	var ipStart, ipEnd net.IP

	if ipStart = net.ParseIP(ipStrStart); ipStart == nil {
		return ErrNotIP
	}
	if ipEnd = net.ParseIP(ipStrEnd); ipEnd == nil {
		return ErrNotIP
	}

	start := binary.BigEndian.Uint32(ipStart.To4())
	end := binary.BigEndian.Uint32(ipEnd.To4())

	if start > end {
		// return decreasing list if range is flipped.
		for i := start; i >= end; i-- {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			fmt.Println(ip)
		}
	} else {
		for i := start; i <= end; i++ {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			fmt.Println(ip)
		}
	}

	return nil
}
