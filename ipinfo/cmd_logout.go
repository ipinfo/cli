package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

func printHelpLogout() {
	fmt.Printf(
		`Usage: %s logout [<opts>]

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdLogout() error {
	var fHelp bool

	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpLogout()
		return nil
	}

	// delete but don't return an error; just log it.
	if err := deleteToken(); err != nil {
		fmt.Println("not logged in")
		return nil
	}

	fmt.Println("logged out")

	return nil
}
