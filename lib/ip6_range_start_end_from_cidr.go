package lib

import (
	"encoding/binary"
	"net"
)

// IP6RangeStartEndFromCIDR returns the start and end IPs in big endian byte
// order of a CIDR in string form.
func IP6RangeStartEndFromCIDR(cidrStr string) (IP6u128, IP6u128, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return IP6u128{}, IP6u128{}, err
	}

	maskhi := binary.BigEndian.Uint64(ipnet.Mask[:8])
	masklo := binary.BigEndian.Uint64(ipnet.Mask[8:])
	starthi := binary.BigEndian.Uint64(ipnet.IP[:8])
	startlo := binary.BigEndian.Uint64(ipnet.IP[8:])
	endhi := (starthi & maskhi) | (maskhi ^ 0xffffffffffffffff)
	endlo := (startlo & masklo) | (masklo ^ 0xffffffffffffffff)

	start := IP6u128{Hi: starthi, Lo: startlo}
	end := IP6u128{Hi: endhi, Lo: endlo}
	return start, end, nil
}
