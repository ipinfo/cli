package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// cmdToolIsV4In6 is the handler for the "is_v4in6" command.
var completionsToolIs4In6 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// printHelpToolIsV4In6 prints the help message for the "is_v4in6" command.
func printHelpToolIsV4In6() {
	fmt.Printf(
		`Usage: %s tool is_v4in6 [<opts>] <ips>

Description:
  Reports whether given ip is an IPv4-mapped IPv6 address.

Examples:
  %[1]s is_v4in6 "::7f00:1"
  %[1]s is_v4in6 "::ffff:
  %[1]s is_v4in6 "::ffff:
  %[1]s is_v4in6 "aaaa::7f00:1"

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

// cmdToolIsV4In6 is the handler for the "is_v4in6" command.
func cmdToolIsV4In6() error {
	f := lib.CmdToolIsV4In6Flags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsV4In6(f, pflag.Args()[2:], printHelpToolIsV4In6)
}
