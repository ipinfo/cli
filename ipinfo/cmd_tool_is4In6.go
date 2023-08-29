package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// cmdToolIs4In6 is the handler for the "is4In6" command.
var completionsToolIs4In6 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// printHelpToolIs4In6 prints the help message for the "is4In6" command.
func printHelpToolIs4In6() {
	fmt.Printf(
		`Usage: %s tool is4In6 [<opts>] <ips>

Description:
  Reports whether given ip is an IPv4-mapped IPv6 address.

Examples:
  %[1]s is4In6 "::7f00:1"
  %[1]s is4In6 "::ffff:
  %[1]s is4In6 "::ffff:
Options:
  General:
    --help, -h
      show help.
`, progBase)
}

// cmdToolIs4In6 is the handler for the "is4In6" command.
func cmdToolIs4In6() error {
	f := lib.CmdToolIs4In6Flags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIs4In6(f, pflag.Args()[2:], printHelpToolIs4In6)
}
