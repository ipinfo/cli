package iputil

import (
	"strconv"
	"strings"
)

// StrIsASNStr checks whether an ASN string really is an ASN of the form
// "asX" or "ASX" where "X" is the ASN's number.
func StrIsASNStr(asn string) bool {
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
