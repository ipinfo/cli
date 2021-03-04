package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

func printHelpASN(asn string) {
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

  Formats:
    --json, -j
      output JSON format. (default)
`, progBase, asn)
}

func cmdASN(asn string) error {
	var fTok string
	var fHelp bool
	var fJSON bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVarP(&fJSON, "json", "j", true, "output JSON format. (default)")
	pflag.Parse()

	if fHelp {
		printHelpASN(asn)
		return nil
	}

	if err := prepareIpinfoClient(fTok); err != nil {
		return err
	}

	data, err := ii.GetASNDetails(asn)
	if err != nil {
		return err
	}

	return outputJSON(data)
}
