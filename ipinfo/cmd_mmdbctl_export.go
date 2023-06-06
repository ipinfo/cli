package main

import (
	"fmt"

	"github.com/ipinfo/mmdbctl/lib"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var predictFormats = []string{"csv", "tsv", "json"}

var completionsExport = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":          predict.Nothing,
		"--help":      predict.Nothing,
		"-o":          predict.Nothing,
		"--out":       predict.Nothing,
		"-f":          predict.Set(predictFormats),
		"--format":    predict.Set(predictFormats),
		"--no-header": predict.Nothing,
	},
}

func printHelpExport() {
	fmt.Printf(
		`Usage: %s export [<opts>] <mmdb_file> [<out_file>]

Options:
  General:
    --help, -h
      show help.

  Input/Output:
    -o <fname>, --out <fname>
      output file name. (e.g. out.csv)
      default: <out_file> if specified, otherwise stdout.

  Format:
    -f <format>, --format <format>
      the output file format.
      can be "csv", "tsv" or "json".
      default: csv if output file ends in ".csv", tsv if ".tsv",
      json if ".json", otherwise csv.
    --no-header
      don't output the header for file formats that include one, like
      CSV/TSV/JSON.
      default: false.
`, progBase)
}

func cmdExport() error {
	f := lib.CmdExportFlags{}
	f.Init()
	pflag.Parse()
	if pflag.NArg() <= 2 {
		f.Help = true
	}

	return lib.CmdExport(f, pflag.Args()[2:], printHelpExport)
}
