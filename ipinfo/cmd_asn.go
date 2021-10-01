package main

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
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
    --field <field>, -f <field>
      lookup only specific fields in the output.
      field names correspond to JSON keys, e.g. 'registry' or 'allocated'.
      multiple field names must be separated by commas.
    --nocolor
      disable colored output.

  Formats:
    --json, -j
      output JSON format. (default)
`, progBase, asn)
}

func cmdASN(asn string) error {
	var fTok string
	var fHelp bool
	var fField []string
	var fJSON bool
	var fNoColor bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.StringSliceVarP(&fField, "field", "f", nil, "specific field to lookup.")
	pflag.BoolVarP(&fJSON, "json", "j", true, "output JSON format. (default)")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable color output.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpASN(asn)
		return nil
	}

	ii = prepareIpinfoClient(fTok)

	// require token for ASN API.
	if ii.Token == "" {
		return errors.New("ASN lookups require a token; login via `ipinfo login`.")
	}

	data, err := ii.GetASNDetails(asn)
	if err != nil {
		return err
	}

	if len(fField) > 0 {
		d := make(ipinfo.BatchASNDetails, 1)
		d[data.ASN] = data
		return outputFieldBatchASNDetails(d, fField, false, false)
	}

	return outputJSON(data)
}
