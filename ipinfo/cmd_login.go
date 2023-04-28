package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var completionsLogin = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-t":         predict.Nothing,
		"--token":    predict.Nothing,
		"--no-check": predict.Nothing,
		"-h":         predict.Nothing,
		"--help":     predict.Nothing,
	},
}

func printHelpLogin() {
	fmt.Printf(
		`Usage: %s login [<opts>] [<token>]

Options:
  --token <tok>, -t <tok>
    token to login with.
    (this is potentially unsafe; let the CLI prompt you instead).
  --no-check
    disable checking if the token is valid or not.
    default: false.
  --help, -h
    show help.
`, progBase)
}

func cmdLogin() error {
	var fTok string
	var fNoCheck bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to save.")
	pflag.BoolVar(&fNoCheck, "no-check", false, "disable checking if token is valid.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpLogin()
		return nil
	}

	// get args without subcommand.
	args := pflag.Args()[1:]

	// only token arg allowed.
	if len(args) > 1 {
		return errors.New("invalid arguments")
	}

	// allow only flag or arg for token but not both.
	if fTok != "" && len(args) > 0 {
		return errors.New("ambiguous token input source")
	}

	// get token, from flag or command line.
	tok := fTok
	if len(args) > 0 {
		tok = args[0]
	}
	for tok == "" {
		fmt.Printf("Enter token: ")
		tokraw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return err
		}

		tok = string(tokraw[:])

		// exit if we have a token now.
		if tok != "" {
			break
		}

		fmt.Println("please enter a token")
	}

	// check token validity.
	if !fNoCheck {
		fmt.Println("logging in...")
		tokenOk, err := isTokenValid(tok)
		if err != nil {
			return fmt.Errorf("could not confirm if token is valid: %w", err)
		}
		if !tokenOk {
			return errors.New("invalid token")
		}
	}

	// save token to file.
	gConfig.Token = tok
	if err := SaveConfig(gConfig); err != nil {
		return err
	}

	fmt.Println("done")

	return nil
}

type meResponse struct {
	Error string `json:"error"`
}

func isTokenValid(tok string) (bool, error) {
	// make API req for true token validity.
	res, err := http.Get("http://localhost:3000/me?token=" + tok)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// parse response.
	me := &meResponse{}
	if err := json.NewDecoder(res.Body).Decode(me); err != nil {
		return false, err
	}

	// If no errors then me.Error should be empty
	return me.Error == "", nil
}
