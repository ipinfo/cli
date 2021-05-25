package main

import (
	"fmt"
	"net"
	"os"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

func printHelpDefault() {
	fmt.Printf(
		`Usage: %s <cmd> [<opts>] [<args>]

Commands:
  <ip>        look up details for an IP address, e.g. 8.8.8.8.
  <asn>       look up details for an ASN, e.g. AS123 or as123.
  myip        get details for your IP.
  bulk        get details for multiple IPs in bulk.
  summarize   get summarized data for a group of IPs.
  map         open a URL to a map showing the locations of a group of IPs.
  prips       print IP list from CIDR or range.
  grepip      grep for IPs matching criteria from any source.
  login       save an API token session.
  logout      delete your current API token session.
  version     show current version.

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --help, -h
      show help.

  Outputs:
    --field <field>, -f <field>
      lookup only a specific field in the output.
      field names correspond to JSON keys, e.g. 'hostname' or 'company.type'.
    --nocolor
      disable colored output.

  Formats:
    --pretty, -p
      output pretty format.
    --json, -j
      output JSON format.
    --csv, -c
      output CSV format.
`, progBase)
}

func cmdDefault() (err error) {
	var ips []net.IP
	var fTok string
	var fHelp bool
	var fField string
	var fPretty bool
	var fJSON bool
	var fCSV bool
	var fNoColor bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.StringVarP(&fField, "field", "f", "", "specific field to lookup.")
	pflag.BoolVarP(&fPretty, "pretty", "p", true, "output pretty format.")
	pflag.BoolVarP(&fJSON, "json", "j", true, "output JSON format. (default)")
	pflag.BoolVarP(&fCSV, "csv", "c", false, "output CSV format.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpDefault()
		return nil
	}

	args := pflag.Args()
	if len(args) != 0 && args[0] != "-" {
		fmt.Printf("err: \"%s\" is not a command.\n", os.Args[1])
		fmt.Println()
		printHelpDefault()
		return nil
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		printHelpDefault()
		return nil
	}

	ips = lib.IPsFromStdin()
	if len(ips) == 0 {
		fmt.Println("no input ips")
		return nil
	}

	ii = prepareIpinfoClient(fTok)

	// require token for bulk.
	if ii.Token == "" {
		fmt.Println("bulk lookups require a token")
		return nil
	}

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
