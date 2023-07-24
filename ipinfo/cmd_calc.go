package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
	"os"
)

var completionsCalc = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpCalc() {
	fmt.Printf(
		`Usage: %s calc <expression> [<opts>]

calc <expression>
  Evaluate a mathematical expression and print the result.

Example:
  %[1]s calc "2*2828-1"
  %[1]s calc "190.87.89.1*2"
  %[1]s calc "2001:0db8:85a3:0000:0000:8a2e:0370:7334*6"

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase)
}

func cmdCalc() error {
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

	switch {
	case cmd != "":
		res, err = cmdCalcInfix()
	default:
		printHelpCalc()
	}

	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "err: %v\n", err)
		if err != nil {
			return err
		}

		printHelpCalc()
	}

	fmt.Println(res)
	return nil
}
