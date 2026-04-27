package main

import (
	"fmt"
	"os"
)

func warnIfInitConfigFails() {
	if err := InitConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "warn: error in creating config file.")
	}
}

func init() {
	warnIfInitConfigFails()
}
