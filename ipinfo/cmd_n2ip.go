package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

var completionsN2IP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-6":        predict.Set(predictReadFmts),
		"--ipv6":    predict.Set(predictReadFmts),
	},
}

func printHelpN2IP() {
	fmt.Printf(
		`Usage: %s n2ip [<opts>] <expr>

Example:
  %[1]s n2ip "2*2828-1"
  %[1]s n2ip "190.87.89.1*2"
  %[1]s n2ip "2001:0db8:85a3:0000:0000:8a2e:0370:7334*6"

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
    --ipv6, -6
      force conversion to IPv6 address
`, progBase)
}

func cmdN2IP() error {
	var forceIpv6 bool
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable colored output.")
	pflag.BoolVarP(&forceIpv6, "ipv6", "6", false, "force conversion to IPv6 address")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}
	if fHelp {
		printHelpDefault()
		return nil
	}

	cmd := ""
	// Reading input from the command line
	if forceIpv6 && len(os.Args) > 3 {
		cmd = os.Args[3]
	} else if !forceIpv6 && len(os.Args) > 2 {
		cmd = os.Args[2]
	} else {
		printHelpN2IP()
		return nil
	}

	// Validate the input
	if strings.TrimSpace(cmd) == "" {
		printHelpN2IP()
		return nil
	}

	res, err := lib.CmdN2IP(cmd, forceIpv6)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "err: %v\n", err)
		if err != nil {
			return err
		}
	}

	fmt.Println(res)
	return nil
}
