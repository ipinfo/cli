package main

import (
	"fmt"
	"net"

	"github.com/pkg/browser"
	"github.com/spf13/pflag"
	"github.com/ipinfo/cli/lib"
)

func printHelpMap() {
	fmt.Printf(
		`Usage: %s map [<opts>] <paths or '-' or cidrs or ip-range>

Description:
  Accepts file paths, '-' for stdin, CIDRs and IP ranges.

  # Lookup all IPs from stdin ('-' can be implied).
  $ %[1]s prips 8.8.8.0/24 | %[1]s map
  $ %[1]s prips 8.8.8.0/24 | %[1]s map -

  # Lookup all IPs in 2 files.
  $ %[1]s map /path/to/iplist1.txt /path/to/iplist2.txt

  # Lookup all IPs from CIDR.
  $ %[1]s map 8.8.8.0/24

  # Lookup all IPs from multiple CIDRs.
  $ %[1]s map 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # Lookup all IPs in an IP range.
  $ %[1]s map 8.8.8.0 8.8.8.255

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

func cmdMap() (err error) {
	var ips []net.IP
	var fHelp bool

	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpMap()
		return nil
	}

	ips, err = lib.IPsFromAllSources(pflag.Args()[1:])
	if err != nil {
		return err
	}
	if len(ips) == 0 {
		fmt.Println("no input ips")
		return nil
	}

	ii = prepareIpinfoClient("")
	d, err := ii.GetIPMap(ips)
	if err != nil {
		return err
	}
	if err := browser.OpenURL(d.ReportURL); err != nil {
		// if it fails, just print the URL.
		fmt.Println(d.ReportURL)
	}

	return nil
}
