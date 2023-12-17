package iputil

import (
	"net"
)

// StrIsIP6Str checks whether a string is an IPv6.
func StrIsIP6Str(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && len(ip) == net.IPv6len
}
