package lib

import (
	"net"
	"regexp"
)

// StrIsIPStr checks whether a string is an IP.
func StrIsIPStr(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}

// StrIsIPv4Str checks if the given string is an IPv4 address
func StrIsIPv4Str(expression string) bool {
	// Compile the regular expression
	ipV4Regex := regexp.MustCompile(ipV4RgxPattern)

	// Use the MatchString function to check if the expression matches the IPv4 pattern
	return ipV4Regex.MatchString(expression)
}

// StrIsIPv6Str checks if the given string is an IPv6 address
func StrIsIPv6Str(expression string) bool {
	// Compile the regular expression
	ipV6Regex := regexp.MustCompile(ipV6RgxPattern)

	// Use the MatchString function to check if the expression matches the IPv6 pattern
	return ipV6Regex.MatchString(expression)
}
