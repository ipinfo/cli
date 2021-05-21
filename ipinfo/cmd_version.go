package main

import (
	"fmt"

	"github.com/ipinfo/complete/v3"
)

var completionsVersion = &complete.Command{}

func cmdVersion() error {
	fmt.Println(version)
	return nil
}
