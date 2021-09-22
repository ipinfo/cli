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

	tokenOk, err := isTokenReal(tok)
	if err != nil {
		return err
	}

	if !tokenOk {
		return errors.New("invalid token")
	}

	// save token to file.
	if err := saveToken(tok); err != nil {
		return err
	}

	fmt.Println("logged in")

	return nil
}

// Custom struct for the response of /me
type MeResponse struct {
	Error string `json:"error"`
	Token string `json:"token"`
}

func isTokenReal(tok string) (bool, error) {
	// Make a request to the /me ep
	res, err := http.Get("https://ipinfo.io/me?token=" + tok)
	if err != nil {
		return false, err
	}

	// Read the body of the response
	me := &MeResponse{}

	if err := json.NewDecoder(res.Body).Decode(me); err != nil {
		return false, err
	}

	// If no errors then me.Error should be empty
	return len(me.Error) == 0, nil
}
