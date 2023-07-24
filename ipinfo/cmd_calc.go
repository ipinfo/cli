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
  %s calc "2*2828-1"
  %s calc "190.87.89.1*2"
  %s calc "2001:0db8:85a3:0000:0000:8a2e:0370:7334*6"

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase, progBase, progBase, progBase)
}

func calcHelp() (err error) {
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

	// currently we do nothing by default.
	printHelpCalc()
	return nil
}

func cmdCalc() error {
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
		err = calcHelp()
	}

	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "err: %v\n", err)
		if err != nil {
			return err
		}
	}

	fmt.Println(res)

	return nil
}
