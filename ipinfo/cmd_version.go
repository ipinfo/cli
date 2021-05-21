package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib/complete"
)

var completionsVersion = &complete.Command{}

func cmdVersion() error {
	fmt.Println(version)
	return nil
}
