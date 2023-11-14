package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsGrepDom = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-o":              predict.Nothing,
		"--only-matching": predict.Nothing,
		"-h":              predict.Nothing,
		"--no-filename":   predict.Nothing,
		"--no-recurse":    predict.Nothing,
		"--help":          predict.Nothing,
		"--nocolor":       predict.Nothing,
		"--punycode":      predict.Nothing,
	},
}

func printHelpGrepDom() {
	fmt.Printf(
		`Usage: %s grepdom [<opts>]

Options:
  General:
    --only-matching, -o
      print only matched domains in result line, excluding surrounding content.
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
  --no-punycode, -n
  	do not convert domains to punycode.
`, progBase)
}

func cmdGrepDom() error {
	f := lib.CmdGrepDomainFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdGrepDomain(f, pflag.Args()[1:], printHelpGrepDom)
}
