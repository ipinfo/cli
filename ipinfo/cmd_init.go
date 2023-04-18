package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"golang.org/x/term"
)

func printHelpInit() {
	fmt.Printf(
		`Usage: %s init [<opts>] [<token>]

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

func cmdInit() error {
	var num int
	var fTok string
	var fNoCheck bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to save.")
	pflag.BoolVar(&fNoCheck, "no-check", false, "disable checking if token is valid.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpInit()
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
	if tok != "" && !fNoCheck {
		if err := checkValidity(tok); err != nil {
			return fmt.Errorf("could not confirm if token is valid: %w", err)
		} else {
			return nil
		}
	}
	fmt.Printf("1) Enter an existing API token\n")
	fmt.Printf("2) Create a new account\n")
	_, err := fmt.Scanf("%d", &num)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return err
	}
	if num == 1 {
		if tok, err = enterToken(tok); err != nil {
			return fmt.Errorf(err.Error())
		}
		if err := checkValidity(tok); err != nil {
			return fmt.Errorf("could not confirm if token is valid: %w", err)
		}
	} else if num == 2 {
		fmt.Printf("Signup flow")
	} else {
		fmt.Println("Invalid input.")
		return err
	}

	return nil
}

func checkValidity(tok string) error {
	fmt.Println("logging in...")
	tokenOk, err := isTokenValid(tok)
	if err != nil {
		return fmt.Errorf("could not confirm if token is valid: %w", err)
	}
	if !tokenOk {
		return errors.New("invalid token")
	}

	// save token to file.
	gConfig.Token = tok
	if err := SaveConfig(gConfig); err != nil {
		return err
	}

	fmt.Println("done")

	return nil
}

func enterToken(tok string) (string, error) {
	for tok == "" {
		fmt.Printf("Enter token: ")
		tokraw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}

		tok = string(tokraw[:])

		// exit if we have a token now.
		if tok != "" {
			break
		}

		fmt.Println("please enter a token")
	}

	return tok , nil
}
