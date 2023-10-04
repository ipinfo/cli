package main

import (
	"fmt"
	"os"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolPrefix = &complete.Command{
	Sub: map[string]*complete.Command{
		"addr":     completionsToolPrefixAddr,
		"bits":     completionsToolPrefixBits,
		"masked":   completionsToolPrefixMasked,
		"is_valid": completionsToolPrefixIsValid,
	},
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpToolPrefix() {
	fmt.Printf(
		`Usage: %s tool prefix <cmd> [<opts>] [<args>]

Commands:
  addr      returns the base IP address of a prefix.
  bits      returns the length of a prefix and reports -1 if invalid.
  masked    returns canonical form of a prefix, masking off non-high bits, and returns the zero if invalid.
  is_valid  reports whether a prefix is valid.

Options:
  --help, -h
    show help.
`, progBase)
}

func toolPrefixHelp() (err error) {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpToolPrefix()
		return nil
	}

	printHelpToolPrefix()
	return nil
}

func cmdToolPrefix() error {
	var err error
	cmd := ""
	if len(os.Args) > 3 {
		cmd = os.Args[3]
	}

	switch {
	case cmd == "addr":
		err = cmdToolPrefixAddr()
	case cmd == "bits":
		err = cmdToolPrefixBits()
	case cmd == "masked":
		err = cmdToolPrefixMasked()
	case cmd == "is_valid":
		err = cmdToolPrefixIsValid()
	default:
		err = toolPrefixHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}

	return nil
}
