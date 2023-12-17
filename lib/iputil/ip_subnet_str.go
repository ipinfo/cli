package iputil

import (
	"net"
)

/*
Note that we always represent subnets as CIDRs, hence the filename having to do
with subnets but the functions using "CIDR" in their names.
*/

// StrIsCIDRStr checks whether a string is in proper CIDR form.
func StrIsCIDRStr(cidrStr string) bool {
	_, _, err := net.ParseCIDR(cidrStr)
	return err == nil
}
