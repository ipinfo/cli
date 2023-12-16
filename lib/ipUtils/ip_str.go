package ipUtils

import (
	"net"
)

// StrIsIPStr checks whether a string is an IP.
func StrIsIPStr(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}

// StrIsIPv4Str checks if the given string is an IPv4 address
func StrIsIPv4Str(expression string) bool {
	return IpV4Regex.MatchString(expression)
}

// StrIsIPv6Str checks if the given string is an IPv6 address
func StrIsIPv6Str(expression string) bool {
	return IpV6Regex.MatchString(expression)
}
