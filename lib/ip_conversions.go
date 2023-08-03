package lib

import (
	"fmt"
	"math/big"
	"net"
	"strconv"
)

// IPtoDecimalStr converts an IP address to a decimal string
func IPtoDecimalStr(strIP string) (string, error) {
	if StrIsIPv6Str(strIP) {
		ip := net.ParseIP(strIP)
		if ip == nil {
			return "", ErrNotIP
		}

		decimalIP := IP6toInt(ip)
		return decimalIP.String(), nil
	}
	if StrIsIPv4Str(strIP) {
		ip := net.ParseIP(strIP)
		if ip == nil {
			return "", ErrNotIP
		}
		return strconv.FormatInt(IP4toInt(ip), 10), nil
	} else {
		return "", ErrInvalidInput
	}
}

// DecimalStrToIP converts a decimal string to an IP address
func DecimalStrToIP(decimal string, forceIPv6 bool) (net.IP, error) {
	// Create a new big.Int with a value of 'decimal'
	num := new(big.Int)
	num, success := num.SetString(decimal, 10)

	if !success {
		fmt.Print(decimal)
		return nil, ErrInvalidInput
	}
	// Convert to IPv4 if not forcing IPv6 and 'num' is within the IPv4 range
	if !forceIPv6 && num.Cmp(big.NewInt(4294967295)) <= 0 {
		ip := make(net.IP, 4)
		b := num.Bytes()
		copy(ip[4-len(b):], b)
		return ip, nil
	}
	// Convert to IPv6 if 'num' is within the IPv6 range
	maxIpv6 := new(big.Int)
	maxIpv6.SetString("340282366920938463463374607431768211455", 10)
	if num.Cmp(maxIpv6) <= 0 {
		ip := make(net.IP, 16)
		b := num.Bytes()
		copy(ip[16-len(b):], b)
		return ip, nil
	}
	return nil, ErrInvalidInput
}

// IP6toInt converts an IPv6 address to a big.Int
func IP6toInt(IPv6Address net.IP) *big.Int {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(IPv6Address)
	return IPv6Int
}

// IP4toInt converts an IPv4 address to a big.Int
func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}
