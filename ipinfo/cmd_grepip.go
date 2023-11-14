package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsGrepIP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-o":                 predict.Nothing,
		"--only-matching":    predict.Nothing,
		"-c":                 predict.Nothing,
		"--include-cidrs":    predict.Nothing,
		"-r":                 predict.Nothing,
		"--include-ranges":   predict.Nothing,
		"--cidrs-only":       predict.Nothing,
		"--ranges-only":      predict.Nothing,
		"-h":                 predict.Nothing,
		"--no-filename":      predict.Nothing,
		"--no-recurse":       predict.Nothing,
		"--help":             predict.Nothing,
		"--nocolor":          predict.Nothing,
		"-4":                 predict.Nothing,
		"--ipv4":             predict.Nothing,
		"-6":                 predict.Nothing,
		"--ipv6":             predict.Nothing,
		"-x":                 predict.Nothing,
		"--exclude-reserved": predict.Nothing,
	},
}

func printHelpGrepIP() {
	fmt.Printf(
		`Usage: %s grepip [<opts>]

Options:
  General:
    --only-matching, -o
      print only matched IP in result line, excluding surrounding content.
    --include-cidrs, -c
      prints the CIDRs too.
    --include-ranges, -r
      prints the Ranges too.
    --cidrs-only
      prints the CIDRs only.
    --ranges-only
      prints the Ranges only.
    --no-filename, -h
      don't print source of match in result lines when more than 1 source.
    --no-recurse
      don't recurse into more directories in directory sources.
    --help
      show help.

  Outputs:
    --nocolor
      disable colored output.

  Filters:
    --ipv4, -4
      match only IPv4 addresses.
    --ipv6, -6
      match only IPv6 addresses.
    --exclude-reserved, -x
      exclude reserved/bogon IPs.
      full list can be found at https://ipinfo.io/bogon.
`, progBase)
}

func cmdGrepIP() error {
	f := lib.CmdGrepIPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdGrepIP(f, pflag.Args()[1:], printHelpGrepIP)
}
