package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/ipinfo/complete/v3"
	"github.com/ipinfo/complete/v3/predict"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

var completionsMyIP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-t":        predict.Nothing,
		"--token":   predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-f":        predict.Set(coreFields),
		"--field":   predict.Set(coreFields),
		"--nocolor": predict.Nothing,
		"-p":        predict.Nothing,
		"--pretty":  predict.Nothing,
		"-j":        predict.Nothing,
		"--json":    predict.Nothing,
		"-c":        predict.Nothing,
		"--csv":     predict.Nothing,
	},
}

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
    --field <field>, -f <field>
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
`, progBase)
}

func cmdMyIP() error {
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
		printHelpMyIP()
		return nil
	}

	ii = prepareIpinfoClient(fTok)
	data, err := ii.GetIPInfo(nil)
	if err != nil {
		return err
	}

	if fField != "" {
		d := make(ipinfo.BatchCore, 1)
		d[data.IP.String()] = data
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
