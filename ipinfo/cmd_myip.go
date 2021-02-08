package main

import (
	"fmt"

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
	var fPretty bool
	var fJSON bool
	var fCSV bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
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

	if fJSON {
		return outputJSON(data)
	}
	if fCSV {
		return outputCSV(data)
	}

	outputFriendlyCore(data)
	return nil
}
