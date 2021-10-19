package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsCache = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
	Args: predict.Set([]string{
		"clear",
	}),
}

func printHelpCache() {
	fmt.Printf(
		`Usage: %s cache [<opts>] [clear]

Description:
  Manage the local cache that stores results previously seen.

Examples:
  # Clear all data currently in the cache.
  $ %[1]s cache clear

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdCache() error {
	var fHelp bool

	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	args := pflag.Args()[1:]
	if fHelp || len(args) != 1 {
		printHelpCache()
		return nil
	}

	switch strings.ToLower(args[0]) {
	case "clear":
		path, err := BoltdbCachePath()
		if err != nil {
			return fmt.Errorf("issue getting cache db path: %w", err)
		}

		// simply delete the whole thing.
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("error clearing cache: %w", err)
		}
	default:
		fmt.Printf("err: %s is not a valid subcommand\n\n", args[0])
		printHelpCache()
		return nil
	}

	fmt.Println("cache cleared")

	return nil
}
