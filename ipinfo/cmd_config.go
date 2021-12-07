package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

func printHelpConfig() {
	fmt.Printf(
		`Usage: %s config [<key>=<value>...]

Description:
  Change the configurations.

Examples:
  $ %[1]s config cache=disable
  $ %[1]s config token=testtoken cahce=enable

Options:
  --help, -h
    show help.

Configurations:
  cache=<enable | disable>
    Control whether the cache is enabled or disabled.
  token=<tok>
    Save a token for use when querying the API.
    (Token will not be validated).
`, progBase)
}

func cmdConfig() error {
	var fHelp bool

	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	// get args for config and parsing it.
	args := pflag.Args()[1:]
	if fHelp || len(args) < 1 {
		printHelpConfig()
		return nil
	}
	for _, arg := range args {
		configStr := strings.Split(arg, "=")
		if len(configStr) != 2 {
			if configStr[0] == "cache" || configStr[0] == "token" {
				return fmt.Errorf("err: no value provided for key %s", configStr[0])
			}
			return fmt.Errorf("err: invalid key argument %s", configStr[0])
		}
		switch strings.ToLower(configStr[0]) {
		case "cache":
			switch strings.ToLower(configStr[1]) {
			case "enable":
				gConfig.CacheEnabled = true
			case "disable":
				gConfig.CacheEnabled = false
			default:
				return fmt.Errorf("err: invalid value %s for key cache", configStr[1])
			}
		case "token":
			gConfig.Token = configStr[1]
		default:
			return fmt.Errorf("err: invalid key argument %s", configStr[0])
		}
	}
	if err := SaveConfig(gConfig); err != nil {
		return err
	}
	return nil
}
