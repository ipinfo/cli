package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

func printHelpConfig() {
	fmt.Printf(
		`Usage: %s config [<opts>=<enable | disable>]

Options:
    cache=<enable| disable>
    enable or disable cache in config file
  --help, -h
    show help.
Examples:
    %[1]s config cache=disable
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

	if len(pflag.Args()) < 2 {
		printHelpConfig()
		return nil
	}

	// get arg for config and parsing it.
	arg := pflag.Arg(1)
	configStr := strings.Split(arg, "=")
	if len(configStr) != 2 {
		if configStr[0] == "cache" {
			fmt.Printf("err: no value provided for %s\n\n", configStr[0])
			printHelpConfig()
			return nil
		}
		fmt.Printf("err: invalid argument %s\n\n", configStr[0])
		printHelpConfig()
		return nil
	}
	switch strings.ToLower(configStr[0]) {
	case "cache":
		switch strings.ToLower(configStr[1]) {
		case "enable":
			gConfig.Cache = true
			err := SetConfig(gConfig)
			if err != nil {
				return err
			}
		case "disable":
			gConfig.Cache = false
			err := SetConfig(gConfig)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("err: %s invalid value for %s\n\n", configStr[1], configStr[0])
			printHelpConfig()
			return nil
		}
	default:
		fmt.Printf("err: invalid argument %s\n\n", configStr[0])
		printHelpConfig()
		return nil
	}
	return nil
}
