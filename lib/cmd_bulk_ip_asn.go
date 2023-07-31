package lib

import (
	"github.com/spf13/pflag"
)

// CmdBulkIpAsnFlags are flags expected by CmdBulkIpAsn.
type CmdBulkIpAsnFlags struct {
	token string
}

// Init initializes the common flags available to CmdBulkIpAsn with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdBulkIpAsnFlags) Init() {
	_h := "see description in --help"
	pflag.StringVarP(
		&f.token,
		"token", "t", "",
		_h,
	)
}

func CmdBulkIpAsn(f CmdBulkIpAsnFlags, args []string) error {
	// 1) Separate the args into two groups
	//    a) args that are ip addresses
	//    b) args that are ASNs
	ii = main.prepareIpinfoClient(fTok)
	ips := li.GetIPInfoBatch(args)

	// 2) For each ip address, do a lookup and print the results
	// 3) For each ASN, do a lookup and print the results
	return nil
}
