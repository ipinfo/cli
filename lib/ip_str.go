package lib

import (
	"net"
)

// StrIsIPStr checks whether a string is an IP.
func StrIsIPStr(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}
