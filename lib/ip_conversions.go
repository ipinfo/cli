package lib

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
)

func CalcIP2n(strIP string) (string, error) {
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
