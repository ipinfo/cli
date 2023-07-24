package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsIP2n = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
	},
}

func printHelpIp2n() {
	fmt.Printf(
		`Usage: %s ip2n <ip>

Example:
  %[1]s ip2n "190.87.89.1"
  %[1]s ip2n "2001:0db8:85a3:0000:0000:8a2e:0370:7334
  %[1]s ip2n "2001:0db8:85a3::8a2e:0370:7334
  %[1]s ip2n "::7334
  %[1]s ip2n "7334::""
	

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase)
}

func cmdIP2n() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable colored output.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpDefault()
		return nil
	}

	var err error
	var res string
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	if strings.TrimSpace(cmd) == "" {
		printHelpIp2n()
		return nil
	}

	res, err = lib.CalcIP2n(cmd)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		printHelpIp2n()
		return nil
	}

	fmt.Println(res)
	return nil
}
