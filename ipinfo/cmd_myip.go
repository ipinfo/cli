package main

import (
	"fmt"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

func printHelpMyIP() {
	fmt.Printf(
		`Usage: %s myip [<opts>]

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
    --pretty, -p
      output pretty format. (default)
    --json, -j
      output JSON format.
    --csv, -c
      output CSV format.
`, progBase)
}

func cmdMyIP() error {
	var fTok string
	var fHelp bool
	var fField string
	var fPretty bool
	var fJSON bool
	var fCSV bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.StringVarP(&fField, "field", "f", "", "specific field to lookup.")
	pflag.BoolVarP(&fPretty, "pretty", "p", true, "output pretty format.")
	pflag.BoolVarP(&fJSON, "json", "j", false, "output JSON format.")
	pflag.BoolVarP(&fCSV, "csv", "c", false, "output CSV format.")
	pflag.Parse()

	if fHelp {
		printHelpMyIP()
		return nil
	}

	if err := prepareIpinfoClient(fTok); err != nil {
		return err
	}

	data, err := ii.GetIPInfo(nil)
	if err != nil {
		return err
	}

	if fField != "" {
		d := make(ipinfo.BatchCore, 1)
		d[data.IP.String()] = data
		return outputFieldBatchCore(d, fField)
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
