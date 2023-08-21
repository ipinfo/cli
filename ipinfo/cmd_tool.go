package main

import (
	"fmt"
	"os"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsTool = &complete.Command{
	Sub: map[string]*complete.Command{
		"aggregate": completionsToolAggregate,
		"is_v4":     completionsToolIs_v4,
		"is_v6":     completionsToolIs_v6,
	},
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpTool() {

	fmt.Printf(
		`Usage: %s tool <cmd> [<opts>] [<args>]

Commands:
  aggregate    aggregate IPs, IP ranges, and CIDRs.
  is_v4        reports whether ip is an IPv4 address.
  is_v6        reports whether ip is an IPv6 address.

Options:
  --help, -h
    show help.
`, progBase)
}

func toolHelp() (err error) {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpTool()
		return nil
	}

	printHelpTool()
	return nil
}

func cmdTool() error {
	var err error
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	switch {
	case cmd == "aggregate":
		err = cmdToolAggregate()
	case cmd == "is_v4":
		err = cmdToolIs_v4()
	case cmd == "is_v6":
		err = cmdToolIs_v6()
	default:
		err = toolHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}

	return nil
}
