package main

import (
	"fmt"
	"net"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
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
    --field, -f
      lookup only a specific field in the output.
      field names correspond to JSON keys, e.g. 'hostname' or 'company.type'.
    --nocolor
      disable colored output.

  Formats:
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
	pflag.BoolVarP(&fJSON, "json", "j", false, "output JSON format.")
	pflag.BoolVarP(&fCSV, "csv", "c", false, "output CSV format.")
	pflag.BoolVarP(&fNoColor, "nocolor", "", false, "disable color output.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpIP(ipStr)
		return nil
	}

	ip := net.ParseIP(ipStr)
	ii = prepareIpinfoClient(fTok)
	data, err := ii.GetIPInfo(ip)
	if err != nil {
		return err
	}

	if fField != "" {
		d := make(ipinfo.BatchCore, 1)
		d[ipStr] = data
		return outputFieldBatchCore(d, fField, false, true)
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
