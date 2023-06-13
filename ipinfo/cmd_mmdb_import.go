package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/ipinfo/mmdbctl/lib"
	"github.com/spf13/pflag"
)

var predictIpVsn = []string{"4", "6"}
var predictSize = []string{"24", "28", "32"}
var predictMerge = []string{"none", "toplevel", "recurse"}

var completionsMmdbImport = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":                          predict.Nothing,
		"--help":                      predict.Nothing,
		"-i":                          predict.Nothing,
		"--in":                        predict.Nothing,
		"-o":                          predict.Nothing,
		"--out":                       predict.Nothing,
		"-c":                          predict.Nothing,
		"--csv":                       predict.Nothing,
		"-t":                          predict.Nothing,
		"--tsv":                       predict.Nothing,
		"-j":                          predict.Nothing,
		"--json":                      predict.Nothing,
		"-f":                          predict.Nothing,
		"--fields":                    predict.Nothing,
		"--fields-from-header":        predict.Nothing,
		"--range-multicol":            predict.Nothing,
		"--joinkey-col":               predict.Nothing,
		"--no-fields":                 predict.Nothing,
		"--no-network":                predict.Nothing,
		"--ip":                        predict.Set(predictIpVsn),
		"-s":                          predict.Set(predictSize),
		"--size":                      predict.Set(predictSize),
		"-m":                          predict.Set(predictMerge),
		"--merge":                     predict.Set(predictMerge),
		"--ignore-empty-values":       predict.Nothing,
		"--disallow-reserved":         predict.Nothing,
		"--alias-6to4":                predict.Nothing,
		"--disable-metadata-pointers": predict.Nothing,
	},
}

func printHelpMmdbImport() {
	fmt.Printf(
		`Usage: %s mmdb import [<opts>] [<input>] [<output>]

Example:
  # Imports an input file and outputs an mmdb file with default configurations. 
  $ %[1]s mmdb import input.csv output.mmdb

Options:
  General:
    --help, -h
      show help.

  Input/Output:
    -i <fname>, --in <fname>
      input file name. (e.g. data.csv or - for stdin)
      must be in CSV, TSV or JSON.
      default: stdin.
    -o <fname>, --out <fname>
      output file name. (e.g. sample.mmdb)
      default: stdout.
    -c, --csv
      interpret input file as CSV.
      by default, the .csv extension will turn this on.
    -t, --tsv
      interpret input file as TSV.
      by default, the .tsv extension will turn this on.
    -j, --json
      interpret input file as JSON.
      by default, the .json extension will turn this on.

  Fields:
    One of the following fields flags, or other flags that implicitly specify
    these, must be used, otherwise --fields-from-header is assumed.

    The first field is always implicitly the network field, unless
    --range-multicol is used, in which case the first 2 fields are considered
    to be start_ip,end_ip.

    When specifying --fields, do not specify the network column(s).

    -f, --fields <comma-separated-fields>
      explicitly specify the fields to assume exist in the input file.
      example: col1,col2,col3
      default: N/A.
    --fields-from-header
      assume first line of input file is a header, and set the fields as that.
      default: true if no other field source is used, false otherwise.
    --range-multicol
      assume that the network field is actually two columns start_ip,end_ip.
      default: false.
    --joinkey-col
      assume --range-multicol and that the 3rd column is join_key, and ignore
      this column when converting to JSON.
      default: false.
    --no-fields
      specify that no fields exist except the implicit network field.
      when enabled, --no-network has no effect; the network field is written.
      default: false.
    --no-network
      if --fields-from-header is set, then don't write the network field, which
      is assumed to be the *first* field in the header.
      default: false.

  Meta:
    --ip <4 | 6>
      output file's ip version.
      default: 6.
    -s, --size <24 | 28 | 32>
      size of records in the mmdb tree.
      default: 32.
    -m, --merge <none | toplevel | recurse>
      the merge strategy to use when inserting entries that conflict.
        none     => no merge; only replace conflicts.
        toplevel => merge only top-level keys.
        recurse  => recursively merge.
      default: none.
    --ignore-empty-values
      if enabled, write into /0 with empty values for all fields, and for any
      entry, don't write out a field whose value is the empty string.
      default: false.
    --disallow-reserved
      disallow reserved networks to be added to the tree.
      default: false.
    --alias-6to4
      enable the mapping of some IPv6 networks into the IPv4 network, e.g.
      ::ffff:0:0/96, 2001::/32 & 2002::/16.
      default: false.
    --disable-metadata-pointers
      some mmdb readers fail to properly read pointers within metadata. this
      allows turning off such pointers.
      NOTE: on by default until we use a different reader in the data repo.
      default: true.
`, progBase)
}

func cmdMmdbImport() error {
	f := lib.CmdImportFlags{}
	f.Init()
	pflag.Parse()
	if pflag.NArg() <= 2 && pflag.NFlag() == 0 {
		f.Help = true
	}

	return lib.CmdImport(f, pflag.Args()[2:], printHelpMmdbImport)
}
