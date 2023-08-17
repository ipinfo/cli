package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsASNBulk = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-t":        predict.Nothing,
		"--token":   predict.Nothing,
		"--nocache": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-f":        predict.Set(asnFields),
		"--field":   predict.Set(asnFields),
		"-j":        predict.Nothing,
		"--json":    predict.Nothing,
	},
}

// printHelpASNBulk prints the help message for the asn bulk command.
func printHelpASNBulk() {
	fmt.Printf(
		`Usage: %s asn bulk [<opts>] <ASNs | filepath>

Description:
  Accepts ASNs and file paths.

Examples:
  # Lookup all ASNs in multiple files.
  $ %[1]s asn bulk /path/to/asnlist1.txt /path/to/asnlist2.txt

  # Lookup multiple ASNs.
  $ %[1]s asn bulk AS123 AS456 AS789

  # Lookup ASNs from multiple sources simultaneously.
  $ %[1]s asn bulk AS123 AS456 AS789 asns.txt

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
      field names correspond to JSON keys, e.g. 'name' or 'registry'.
      multiple field names must be separated by commas.

  Formats:
    --json, -j
      output JSON format. (default)
    --yaml, -y
      output YAML format.
`, progBase)
}

// cmdASNBulk is the asn bulk command.
func cmdASNBulk(piped bool) error {
	f := lib.CmdASNBulkFlags{}
	f.Init()
	pflag.Parse()

	ii = prepareIpinfoClient(f.Token)
	var args []string
	if !piped {
		args = pflag.Args()[2:]
	}

	data, err := lib.CmdASNBulk(f, ii, args, printHelpASNBulk)
	if err != nil {
		return err
	}
	if (data) == nil {
		return nil
	}

	if len(f.Field) > 0 {
		return outputFieldBatchASNDetails(data, f.Field, false, false)
	}

	if f.Yaml {
		return outputYAML(data)
	}

	return outputJSON(data)
}
