package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

var completionsN2IP6 = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-f":        predict.Set(predictReadFmts),
		"--format":  predict.Set(predictReadFmts),
	},
}

func printHelpN2IP6() {

	fmt.Printf(
		`Usage: %s n2ip6 [<opts>] <expr>

Example:
  %s n2ip6 "190.87.89.1"
  %s n2ip6 "2001:0db8:85a3:0000:0000:8a2e:0370:7334
  %s n2ip6 "2001:0db8:85a3::8a2e:0370:7334
  %s n2ip6 "::7334
  %s n2ip6 "7334::""
	

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase, progBase, progBase, progBase, progBase, progBase)
}

func n2ip6Help() (err error) {
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
	printHelpN2IP6()
	return nil
}

func cmdN2IP6() error {
	var err error

	cmd := ""

	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	if strings.TrimSpace(cmd) == "" {
		err := n2ip6Help()
		if err != nil {
			return err
		}
		return nil
	}

	if isInvalid(cmd) {
		return errors.New("invalid expression")
	}
	tokens, err := tokeinzeExp(cmd)

	if err != nil {
		return err
	}

	postfix := infixToPostfix(tokens)

	result, err := evaluatePostfix(postfix)

	if err != nil {
		return err
	}

	res := decimalToIP(result.String(), true)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		err := n2ip6Help()
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Println(res)

	return nil
}
