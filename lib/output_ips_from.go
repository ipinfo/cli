package lib

import (
	"encoding/binary"
	"fmt"
	"net"
)

// OutputIPsFrom outputs a list of IPs from inputs which are interpreted to
// contain IP ranges and IP CIDRs in them, all depending upon which flags are
// set.
func OutputIPsFrom(
	inputs []string,
	iprange bool,
	cidr bool,
) error {
	// prevent edge cases with all flags turned off.
	if !iprange && !cidr {
		return nil
	}

	for _, input := range inputs {
		if iprange {
			if err := OutputIPsFromRangeStr(input); err == nil {
				continue
			}
		}

		if cidr && StrIsCIDR(input) {
			if err := OutputIPsFromCIDR(input); err == nil {
				continue
			}
		}

		return ErrInvalidInput
	}

	return nil
}

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

// OutputIPsFromRange is the same as IPsFromRange with O(1) memory by
// discarding IPs after printing.
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

// OutputIPsFromRangeStr outputs all IPs in an IP range.
// The string must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func OutputIPsFromRangeStr(r string) error {
	rStart, rEnd, err := IPRangeStrFromRangeStr(r)
	if err != nil {
		return err
	}

	return OutputIPsFromRange(rStart, rEnd)
}
