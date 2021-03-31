package lib

import (
	"encoding/binary"
	"net"
)

// IPRangeStartEndFromCIDR returns the start and end IPs in big endian byte
// order of a CIDR in string form.
func IPRangeStartEndFromCIDR(cidrStr string) (uint32, uint32, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return 0, 0, err
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)
	end := (start & mask) | (mask ^ 0xffffffff)
	return start, end, nil
}
