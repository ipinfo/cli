package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsSingleIp = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsSingleIp() {
	fmt.Printf(
		`Usage: %s tool isSingleIp [<opts>] <cidr | filepath>

Description:
  checks whether a CIDR contains exactly one IP.

Examples:
  # Check CIDR.
  $ %[1]s tool isSingleIp 1.1.1.0/30

  # Check for file.
  $ %[1]s tool isSingleIp /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool isSingleIp

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdToolIsSingleIp() (err error) {
	f := lib.CmdToolIsSingleIpFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsSingleIp(f, pflag.Args()[2:], printHelpToolIsSingleIp)
}
