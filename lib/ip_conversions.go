package lib

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"regexp"
	"strconv"
)

func IPtoDecimal(strIP string) (string, error) {
	if IsIPv6Address(strIP) {
		ip := net.ParseIP(strIP)
		if ip == nil {
			return "", errors.New("invalid IPv6 address: '" + strIP + "'")
		}

		decimalIP := IP6toInt(ip)
		return decimalIP.String(), nil
	}
	if IsIPv4Address(strIP) {
		ip := net.ParseIP(strIP)
		if ip == nil {
			return "", errors.New("invalid IPv4 address: '" + strIP + "'")
		}
		return strconv.FormatInt(IP4toInt(ip), 10), nil
	} else {
		return "", errors.New("invalid IP address: '" + strIP + "'")
	}
}

func DecimalToIP(decimal string, forceIPv6 bool) net.IP {
	// Create a new big.Int with a value of 'decimal'
	num := new(big.Int)
	num, success := num.SetString(decimal, 10)
	if !success {
		fmt.Fprintf(os.Stderr, "Error parsing the decimal string: %v\n", success)
		return nil
	}

	// Convert to IPv4 if not forcing IPv6 and 'num' is within the IPv4 range
	if !forceIPv6 && num.Cmp(big.NewInt(4294967295)) <= 0 {
		ip := make(net.IP, 4)
		b := num.Bytes()
		copy(ip[4-len(b):], b)
		return ip
	}

	// Convert to IPv6
	ip := make(net.IP, 16)
	b := num.Bytes()
	copy(ip[16-len(b):], b)
	return ip
}

func IP6toInt(IPv6Address net.IP) *big.Int {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(IPv6Address)
	return IPv6Int
}
func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}

func IsIPv4Address(expression string) bool {
	// Define the regular expression pattern for matching IPv4 addresses
	ipV4Pattern := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	// Compile the regular expression
	ipV4Regex := regexp.MustCompile(ipV4Pattern)

	// Use the MatchString function to check if the expression matches the IPv4 pattern
	return ipV4Regex.MatchString(expression)
}

func IsIPv6Address(expression string) bool {
	// Define the regular expression pattern for matching IPv6 addresses
	ipV6Pattern := `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`

	// Compile the regular expression
	ipV6Regex := regexp.MustCompile(ipV6Pattern)

	// Use the MatchString function to check if the expression matches the IPv6 pattern
	return ipV6Regex.MatchString(expression)
}
