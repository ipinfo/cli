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
		"next":      completionsToolNext,
		"prev":      completionsToolPrev,
		"ip2n":      completionsToolIP2n,
		"n2ip":      completionsToolN2IP,
		"n2ip6":     completionsToolN2IP6,
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
  next         get next IP of IPs
  prev         get previous IP of IPs
  ip2n         converts an IPv4 or IPv6 address to its decimal representation.
  n2ip	       evaluates a mathematical expression and converts it to an IPv4 or IPv6.
  n2ip6	       evaluates a mathematical expression and converts it to an IPv6.

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
	case cmd == "next":
		err = cmdToolNext()
	case cmd == "prev":
		err = cmdToolPrev()
	case cmd == "ip2n":
		err = cmdToolIP2n()
	case cmd == "n2ip":
		err = cmdToolN2IP()
	case cmd == "n2ip6":
		err = cmdToolN2IP6()
	default:
		err = toolHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}

	return nil
}
