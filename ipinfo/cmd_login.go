package main

import (
	"fmt"
	"os"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var completionsLogin = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-t":      predict.Nothing,
		"--token": predict.Nothing,
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
	},
}

func printHelpLogin() {
	fmt.Printf(
		`Usage: %s login [<opts>]

Options:
  --token <tok>, -t <tok>
    token to login with.
    (this is potentially unsafe; let the CLI prompt you instead).
  --help, -h
    show help.
`, progBase)
}

func cmdLogin() error {
	var fTok string
	var fHelp bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to save.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpLogin()
		return nil
	}

	// get token, from flag or command line.
	tok := fTok
	if tok == "" {
		fmt.Printf("Enter token: ")
		tokraw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return err
		}

		tok = string(tokraw[:])
	}

	// save token to file.
	if err := saveToken(tok); err != nil {
		return err
	}

	fmt.Println("logged in")

	return nil
}
