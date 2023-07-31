package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

var completionsBulk = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-t":        predict.Nothing,
		"--token":   predict.Nothing,
		"--nocache": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-f":        predict.Set(coreFields),
		"--field":   predict.Set(coreFields),
		"--nocolor": predict.Nothing,
		"-j":        predict.Nothing,
		"--json":    predict.Nothing,
		"-c":        predict.Nothing,
		"--csv":     predict.Nothing,
	},
}

func printHelpBulk() {
	fmt.Printf(
		`Usage: %s bulk [<opts>] <ip | ip-range | cidr | filepath>

Description:
  Accepts IPs, IP ranges, CIDRs and file paths.

Examples:
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
  $ %[1]s bulk 8.8.8.0-8.8.8.255

  # Lookup all IPs from multiple sources simultaneously.
  $ %[1]s bulk 8.8.8.0-8.8.8.255 1.1.1.0/30 123.123.123.123 ips.txt

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
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
    --json, -j
      output JSON format. (default)
    --csv, -c
      output CSV format.
    --yaml, -y
      output YAML format.
`, progBase)
}

func cmdBulk() (err error) {
	var ips []net.IP
	var fTok string
	var fField []string
	var fJSON bool
	var fCSV bool
	var fYAML bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVar(&fNoCache, "nocache", false, "disable the cache.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.StringSliceVarP(&fField, "field", "f", nil, "specific field to lookup.")
	pflag.BoolVarP(&fJSON, "json", "j", true, "output JSON format. (default)")
	pflag.BoolVarP(&fCSV, "csv", "c", false, "output CSV format.")
	pflag.BoolVarP(&fYAML, "yaml", "y", false, "output YAML format.")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable color output.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpBulk()
		return nil
	}

	ips, err = lib.IPListFromAllSrcs(pflag.Args()[1:])
	fmt.Println(ips)
	if err != nil {
		return err
	}
	if len(ips) == 0 {
		fmt.Println("no input ips")
		return nil
	}

	ii = prepareIpinfoClient(fTok)

	// require token for bulk.
	if ii.Token == "" {
		return errors.New("bulk lookups require a token; login via `ipinfo init`.")
	}

	data, err := ii.GetIPInfoBatch(ips, ipinfo.BatchReqOpts{
		TimeoutPerBatch:              60 * 30, // 30min
		ConcurrentBatchRequestsLimit: 20,
	})
	if err != nil {
		return err
	}

	if len(fField) > 0 {
		return outputFieldBatchCore(data, fField, true, true)
	}

	if fCSV {
		return outputCSVBatchCore(data)
	}
	if fYAML {
		return outputYAML(data)
	}

	return outputJSON(data)
}
