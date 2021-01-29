package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func cmdLogout(c *cli.Context) error {
	// delete but don't return an error; just log it.
	if err := deleteToken(); err != nil {
		fmt.Println("not logged in")
		return nil
	}

	fmt.Println("logged out")

	return nil
}
