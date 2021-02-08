package main

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

func printHelpIP(ipStr string) {
	fmt.Printf(
		`Usage: %s %s [<opts>]

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
    --csv, -c
      output CSV format.
`, progBase, ipStr)
}

func cmdIP(ipStr string) error {
	var fTok string
	var fPretty bool
	var fJSON bool
	var fCSV bool
	var fHelp bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fPretty, "pretty", "p", true, "output pretty format.")
	pflag.BoolVarP(&fJSON, "json", "j", false, "output JSON format.")
	pflag.BoolVarP(&fCSV, "csv", "c", false, "output CSV format.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpIP(ipStr)
		return nil
	}

	if err := prepareIpinfoClient(fTok); err != nil {
		return err
	}

	ip := net.ParseIP(ipStr)
	data, err := ii.GetIPInfo(ip)
	if err != nil {
		return err
	}

	if fJSON {
		return outputJSON(data)
	}
	if fCSV {
		return outputCSV(data)
	}

	outputFriendlyCore(data)
	return nil
}
