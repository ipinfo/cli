package main

import (
	"fmt"
	"net"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

func printHelpBulk() {
	fmt.Printf(
		`Usage: %s bulk [<opts>] <paths or '-' or cidrs or ip-range>

Description:
  Accepts file paths, '-' for stdin, CIDRs and IP ranges.

  # Lookup all IPs from stdin ('-' can be implied).
  $ %[1]s prips 8.8.8.0/24 | %[1]s bulk
  $ %[1]s prips 8.8.8.0/24 | %[1]s bulk -

  # Lookup all IPs in 2 files.
  $ %[1]s bulk /path/to/iplist1.txt /path/to/iplist2.txt

  # Lookup all IPs from CIDR.
  $ %[1]s bulk 8.8.8.0/24

  # Lookup all IPs from multiple CIDRs.
  $ %[1]s bulk 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # Lookup all IPs in an IP range.
  $ %[1]s bulk 8.8.8.0 8.8.8.255

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --help, -h
      show help.

  Outputs:
    --field, -f
      lookup only a specific field in the output.
      field names correspond to JSON keys, e.g. 'hostname' or 'company.type'.

  Formats:
    --json, -j
      output JSON format. (default)
    --csv, -c
      output CSV format.
`, progBase)
}

func cmdBulk() (err error) {
	var ips []net.IP
	var fTok string
	var fHelp bool
	var fField string
	var fJSON bool
	var fCSV bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.StringVarP(&fField, "field", "f", "", "specific field to lookup.")
	pflag.BoolVarP(&fJSON, "json", "j", true, "output JSON format. (default)")
	pflag.BoolVarP(&fCSV, "csv", "c", false, "output CSV format.")
	pflag.Parse()

	if fHelp {
		printHelpBulk()
		return nil
	}

	ips, err = getInputIPs(pflag.Args()[1:])
	if err != nil {
		return err
	}
	if len(ips) == 0 {
		fmt.Println("no input ips")
		return nil
	}

	ii = prepareIpinfoClient(fTok)
	data, err := ii.GetIPInfoBatch(ips, ipinfo.BatchReqOpts{})
	if err != nil {
		return err
	}

	if fField != "" {
		return outputFieldBatchCore(data, fField, true, false)
	}

	if fCSV {
		return outputCSVBatchCore(data)
	}

	return outputJSON(data)
}
