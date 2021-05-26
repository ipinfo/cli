package lib

import (
	"net"
)

// IPStrFromIPBE returns the IP string representation of an IP in big-endian
// numerical form.
func IPStrFromIPBE(ip uint32) string {
	return net.IPv4(
		byte(ip>>24),
		byte(ip>>16),
		byte(ip>>8),
		byte(ip),
	).String()
}
