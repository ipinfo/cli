package lib

import (
	"net"
)

// IsIP checks whether a string is an IP.
func IsIP(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}
