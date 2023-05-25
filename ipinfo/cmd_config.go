package main

import (
	"fmt"
	"strings"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsConfig = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpConfig() {
	fmt.Printf(
		`Usage: %s config [<key>=<value>...]

Description:
  Change the configurations.

Examples:
  $ %[1]s config cache=disable
  $ %[1]s config token=testtoken cache=enable

Options:
  --help, -h
    show help.

Configurations:
  cache=<enable | disable>
    Control whether the cache is enabled or disabled.
  open_browser=<enable | disable>
    Control whether the links should open the browser or not.
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
		key := strings.ToLower(configStr[0])
		if len(configStr) != 2 {
			if key == "cache" || key == "token" || key == "open_browser" {
				return fmt.Errorf("err: no value provided for key %s", key)
			}
			return fmt.Errorf("err: invalid key argument %s", key)
		}
		switch key {
		case "cache":
			val := strings.ToLower(configStr[1])
			switch val {
			case "enable":
				gConfig.CacheEnabled = true
			case "disable":
				gConfig.CacheEnabled = false
			default:
				return fmt.Errorf("err: invalid value %s; cache must be 'enabled' or disabled", val)
			}
		case "open_browser":
			val := strings.ToLower(configStr[1])
			switch val {
			case "enable":
				gConfig.OpenBrowser = true
			case "disable":
				gConfig.OpenBrowser = false
			default:
				return fmt.Errorf("err: invalid value %s; open_browser must be 'enabled' or disabled", val)
			}
		case "token":
			gConfig.Token = configStr[1]
		default:
			return fmt.Errorf("err: invalid key argument %s", configStr[0])
		}
	}

	// save config in bulk.
	if err := SaveConfig(gConfig); err != nil {
		return err
	}

	return nil
}
