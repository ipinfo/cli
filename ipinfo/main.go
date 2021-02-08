package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipinfo/go/v2/ipinfo"
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

	cmd := os.Args[1]
	switch {
	case isIP(cmd):
		cmdIP(cmd)
	case isASN(cmd):
		asn := strings.ToUpper(cmd)
		cmdASN(asn)
	case cmd == "myip":
		cmdMyIP()
	case cmd == "bulk":
		cmdBulk()
	case cmd == "sum":
		cmdSum()
	case cmd == "prips":
		cmdPrips()
	case cmd == "login":
		cmdLogin()
	case cmd == "logout":
		cmdLogout()
	case cmd == "v":
		cmdVersion()
	case cmd == "vsn":
		cmdVersion()
	case cmd == "version":
		cmdVersion()
	default:
		fmt.Printf("err: \"%s\" is not a command.\n", cmd)
		printHelp()
	}
}
