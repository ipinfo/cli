package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

var progBase = filepath.Base(os.Args[0])

var ii *ipinfo.Client

func prepareIpinfoClient(tok string) error {
	if tok == "" {
		tok, _ = restoreToken()
	}

	ii = ipinfo.NewClient(nil, nil, tok)
	return nil
}

func printHelp() {
	fmt.Printf(
		`Usage: %s <cmd> [<opts>] [<args>]

Commands:
  <ip>        look up details for an IP address, e.g. 8.8.8.8.
  <asn>       look up details for an ASN, e.g. AS123 or as123.
  myip        get details for your IP.
  bulk        get details for multiple IPs in bulk.
  sum         get summarized data for a group of IPs.
  prips       print IP list from CIDR or range.
  login       save an API token session.
  logout      delete your current API token session.
  version     show current version.

Options:
  --help, -h
    show help.
`, progBase)
}

func main() {
	if len(os.Args) == 1 {
		printHelp()
		return
	}

	var err error
	cmd := os.Args[1]
	switch {
	case isIP(cmd):
		err = cmdIP(cmd)
	case isASN(cmd):
		asn := strings.ToUpper(cmd)
		err = cmdASN(asn)
	case cmd == "myip":
		err = cmdMyIP()
	case cmd == "bulk":
		err = cmdBulk()
	case cmd == "sum":
		err = cmdSum()
	case cmd == "prips":
		err = cmdPrips()
	case cmd == "login":
		err = cmdLogin()
	case cmd == "logout":
		err = cmdLogout()
	case cmd == "version":
		err = cmdVersion()
	default:
		var fHelp bool

		pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
		pflag.Parse()

		if fHelp {
			printHelp()
		} else {
			fmt.Printf("err: \"%s\" is not a command.\n", cmd)
			fmt.Println()
			printHelp()
		}
	}

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
