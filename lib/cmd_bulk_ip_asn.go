package lib

import (
	"fmt"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
	"net"
)

// CmdBulkIpAsnFlags are flags expected by CmdBulkIpAsn.
type CmdBulkIpAsnFlags struct {
	Token string
}

// Init initializes the common flags available to CmdBulkIpAsn with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdBulkIpAsnFlags) Init() {
	_h := "see description in --help"
	pflag.StringVarP(
		&f.Token,
		"token", "t", "",
		_h,
	)
}

func CmdBulkIpAsn(f CmdBulkIpAsnFlags, ii *ipinfo.Client, args []string) error {
	var ips []net.IP
	var asns []string
	var err error

	for _, arg := range args {
		if StrIsIPStr(arg) {
			ips = append(ips, net.ParseIP(arg))
		} else if StrIsASNStr(arg) {
			asns = append(asns, arg)
		}
	}

	fmt.Println(asns)
	_, err = ii.GetASNDetailsBatch(asns, ipinfo.BatchReqOpts{
		TimeoutPerBatch:              60 * 30, // 30min
		ConcurrentBatchRequestsLimit: 20,
	})

	//fmt.Println(data)

	// 2) For each ip address, do a lookup and print the results
	// 3) For each ASN, do a lookup and print the results
	return err
}
