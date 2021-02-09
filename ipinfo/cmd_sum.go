package main

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/pflag"
)

func printHelpSum() {
	fmt.Printf(
		`Usage: %s sum [<opts>] <paths or '-' or cidrs or ip-range>

Description:
  Accepts file paths, '-' for stdin, CIDRs and IP ranges.

  # Lookup all IPs from stdin ('-' can be implied).
  $ %[1]s prips 8.8.8.0/24 | %[1]s sum
  $ %[1]s prips 8.8.8.0/24 | %[1]s sum -

  # Lookup all IPs in 2 files.
  $ %[1]s sum /path/to/iplist1.txt /path/to/iplist2.txt

  # Lookup all IPs from CIDR.
  $ %[1]s sum 8.8.8.0/24

  # Lookup all IPs from multiple CIDRs.
  $ %[1]s sum 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # Lookup all IPs in an IP range.
  $ %[1]s sum 8.8.8.0 8.8.8.255

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --help, -h
      show help.

  Outputs:
    --pretty, -p
      output pretty format. (default)
    --json, -j
      output JSON format.
`, progBase)
}

func cmdSum() (err error) {
	var ips []net.IP
	var fTok string
	var fHelp bool
	var fPretty bool
	var fJSON bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVarP(&fPretty, "pretty", "p", true, "output pretty format. (default)")
	pflag.BoolVarP(&fJSON, "json", "j", false, "output JSON format.")
	pflag.Parse()

	if fHelp {
		printHelpSum()
		return nil
	}

	if err := prepareIpinfoClient(fTok); err != nil {
		return err
	}

	args := pflag.Args()[1:]

	// check for stdin, implied or explicit.
	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Println("** manual input mode **")
			fmt.Println("Enter all IPs, one per line:")
		}
		ips = ipsFromStdin()

		goto lookup
	}

	// check for IP range.
	if isIP(args[0]) {
		if len(args) != 2 {
			return errIPRangeRequiresTwoIPs
		}
		if !isIP(args[1]) {
			return errNotIP
		}

		ips, err = ipsFromRange(args[0], args[1])
		if err != nil {
			return err
		}

		goto lookup
	}

	// check for all CIDRs.
	if isCIDR(args[0]) {
		for _, arg := range args[1:] {
			if !isCIDR(arg) {
				return errNotCIDR
			}
		}

		ips, err = ipsFromCIDRs(args)
		if err != nil {
			return err
		}

		goto lookup
	}

	// check for all filepaths.
	if fileExists(args[0]) {
		for _, arg := range args[1:] {
			if !fileExists(arg) {
				return errNotFile
			}
		}

		ips, err = ipsFromFiles(args)
		if err != nil {
			return err
		}

		goto lookup
	}

lookup:

	if len(ips) == 0 {
		fmt.Println("no input ips")
		return nil
	}

	data, err := ii.GetIPSummary(ips)
	if err != nil {
		return err
	}

	if fJSON {
		return outputJSON(data)
	}

	// TODO pretty
	fmt.Printf("%v\n", data)
	return nil
}
