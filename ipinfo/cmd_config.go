package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

func printHelpConfig() {
	fmt.Printf(
		`Usage: %s config [<opts>] [<enable | disable>]

Options:
  --cache <enable| disable>, -c <enable | disable>
    enable or disable cache in config file
  --help, -h
    show help.
`, progBase)
}

func cmdConfig() error {
	var fCache string
	var fHelp bool

	pflag.StringVarP(&fCache, "cache", "c", "", "cache enable | disable.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpConfig()
		return nil
	}
	switch strings.ToLower(fCache) {
	case "enable":
		gConfig.GlobalCache = true
		err := SetConfig(gConfig)
		if err != nil {
			return err
		}
	case "disable":
		gConfig.GlobalCache = false
		err := SetConfig(gConfig)
		if err != nil {
			return err
		}
	default:
		fmt.Printf("err: invalid argument\n\n")
		printHelpConfig()
		return nil
	}
	return nil
}
