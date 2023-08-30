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
		"is_v4":     completionsToolIsV4,
		"isValid":   completionsToolIsValid,
		"lower":     completionsToolLower,
		"upper":     completionsToolUpper,
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
  next         get the next IP of the input IP
  prev         get the previous IP of the input IP
  is_v4        reports whether input is an IPv4 address.
  is_v6        reports whether input is an IPv6 address.
  isValid      reports whether an IP is valid.
  lower        get start IP of IPs, IP ranges, and CIDRs.
  upper        get end IP of IPs, IP ranges, and CIDRs.
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
	case cmd == "is_v4":
		err = cmdToolIsV4()
	case cmd == "is_v6":
		err = cmdToolIsV6()
	case cmd == "isValid":
		err = cmdToolIsValid()
	case cmd == "lower":
		err = cmdToolLower()
	case cmd == "upper":
		err = cmdToolUpper()
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
