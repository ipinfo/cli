package lib

import (
	"net"
)

// IsCIDR checks whether a string is in proper CIDR form.
func IsCIDR(cidrStr string) bool {
	_, _, err := net.ParseCIDR(cidrStr)
	return err == nil
}
