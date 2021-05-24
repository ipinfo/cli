package lib

import (
	"net"
	"strconv"
	"strings"
)

// StrIsASN checks whether an ASN string really is an ASN of the form "asX" or
// "ASX" where "X" is the ASN's number.
func StrIsASN(asn string) bool {
	// check length.
	if len(asn) < 3 {
		return false
	}

	// ensure "AS" or "as" prefix.
	if !strings.HasPrefix(asn, "AS") && !strings.HasPrefix(asn, "as") {
		return false
	}

	// ensure number suffix.
	asnNumStr := asn[2:]
	if _, err := strconv.Atoi(asnNumStr); err != nil {
		return false
	}

	return true
}

// StrIsCIDR checks whether a string is in proper CIDR form.
func StrIsCIDR(cidrStr string) bool {
	_, _, err := net.ParseCIDR(cidrStr)
	return err == nil
}

// StrIsIP checks whether a string is an IP.
func StrIsIP(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}
