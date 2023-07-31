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
	// Define the regular expression pattern for matching IPv4 addresses
	ipV4Pattern := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	// Compile the regular expression
	ipV4Regex := regexp.MustCompile(ipV4Pattern)

	// Use the MatchString function to check if the expression matches the IPv4 pattern
	return ipV4Regex.MatchString(expression)
}

// StrIsIPv6Str checks if the given string is an IPv6 address
func StrIsIPv6Str(expression string) bool {
	// Define the regular expression pattern for matching IPv6 addresses
	ipV6Pattern := `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`

	// Compile the regular expression
	ipV6Regex := regexp.MustCompile(ipV6Pattern)

	// Use the MatchString function to check if the expression matches the IPv6 pattern
	return ipV6Regex.MatchString(expression)
}
