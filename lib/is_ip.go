package lib

import (
	"net"
)

func IsIP(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}
