package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func cmdLogin(c *cli.Context) error {
	// get token, from flag or command line.
	tok := c.String("token")
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

	return nil
}
