package main

import (
	"fmt"
)

func cmdVersion() error {
	fmt.Println(version)
	return nil
}
