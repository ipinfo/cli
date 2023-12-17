package iputil

import (
	"net"
)

/*
Note that we always represent subnets as CIDRs, hence the filename having to do
with subnets but the functions using "CIDR" in their names.
*/

// StrIsCIDR6Str checks whether a string is in proper CIDR IPv6 form.
func StrIsCIDR6Str(cidrStr string) bool {
	ip, _, err := net.ParseCIDR(cidrStr)
	return err == nil && ip != nil && len(ip) == net.IPv6len
}
