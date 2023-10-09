package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

var completionsMyIP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-t":        predict.Nothing,
		"--token":   predict.Nothing,
		"--nocache": predict.Nothing,
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
		"-6":        predict.Nothing,
		"--ipv6":    predict.Nothing,
	},
}

func printHelpMyIP() {
	fmt.Printf(
		`Usage: %s myip [<opts>]

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --ipv6, -6
      use IPv6 address.
    --nocache
      do not use the cache.
    --help, -h
      show help.

  Outputs:
    --field <field>, -f <field>
      lookup only specific fields in the output.
      field names correspond to JSON keys, e.g. 'hostname' or 'company.type'.
      multiple field names must be separated by commas.
    --nocolor
      disable colored output.

  Formats:
    --pretty, -p
      output pretty format. (default)
    --json, -j
      output JSON format.
    --csv, -c
      output CSV format.
    --yaml, -y
      output YAML format.
`, progBase)
}

func cmdMyIP() error {
	var fTok string
	var fField []string
	var fPretty bool
	var fJSON bool
	var fCSV bool
	var fYAML bool
	var fIPv6 bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVar(&fNoCache, "nocache", true, "disable the cache.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.StringSliceVarP(&fField, "field", "f", nil, "specific field to lookup.")
	pflag.BoolVarP(&fPretty, "pretty", "p", true, "output pretty format.")
	pflag.BoolVarP(&fJSON, "json", "j", false, "output JSON format.")
	pflag.BoolVarP(&fCSV, "csv", "c", false, "output CSV format.")
	pflag.BoolVarP(&fYAML, "yaml", "y", false, "output YAML format.")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable color output.")
	pflag.BoolVarP(&fIPv6, "ipv6", "6", false, "use IPv6 address.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpMyIP()
		return nil
	}

	ii = prepareIpinfoClient(fTok)
	if fIPv6 {
		ii.IPv6 = true
	}

	data, err := ii.GetIPInfo(nil)
	if err != nil {
		return err
	}

	if len(fField) > 0 {
		d := make(ipinfo.BatchCore, 1)
		d[data.IP.String()] = data
		return outputFieldBatchCore(d, fField, false, false)
	}
	if fJSON {
		return outputJSON(data)
	}
	if fCSV {
		return outputCSV(data)
	}
	if fYAML {
		return outputYAML(data)
	}

	outputFriendlyCore(data)
	return nil
}
