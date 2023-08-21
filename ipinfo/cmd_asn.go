package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"os"
)

var completionsASN = &complete.Command{
	Sub: map[string]*complete.Command{
		"bulk": completionsASNBulk,
	},
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpASN() {
	fmt.Printf(
		`Usage: %s asn <cmd> [<opts>] 

Commands:
  bulk        lookup ASNs in bulk

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

func cmdASNDefault() error {
	// check whether the standard input (stdin) is being piped
	// or redirected from another source or whether it's being read from the terminal (interactive mode).
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		printHelpASN()
		return nil
	}
	return cmdASNBulk(true)
}

// cmdASN is the handler for the "asn" command.
func cmdASN() error {
	var err error
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	switch {
	case cmd == "bulk":
		err = cmdASNBulk(false)
	default:
		err = cmdASNDefault()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}

	return nil
}
