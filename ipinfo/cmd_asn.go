package main

import (
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
	"os"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
)

var completionsASN = &complete.Command{
	Sub: map[string]*complete.Command{
		"bulk": completionsASNBulk,
	},
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpASN() {

	fmt.Printf(
		`Usage: %s asn <cmd> [<opts>] 

Commands:
  bulk        lookup ASNs in bulk

Options:
  General:
    --help, -h
      show help.
`, progBase)
}

// cmdASN is the handler for the "asn" command.
func cmdASN() error {
	var err error
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	switch {
	case cmd == "bulk":
		err = cmdASNBulk()
	default:
		// check whether the standard input (stdin) is being piped
		// or redirected from another source or whether it's being read from the terminal (interactive mode).
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			printHelpASN()
			return nil
		}

		f := lib.CmdASNBulkFlags{}
		f.Init()
		pflag.Parse()

		ii = prepareIpinfoClient(f.Token)
		data, err := lib.CmdASNBulk(f, ii, lib.ReadStringsFromStdin(), printHelpASNBulk)
		if err != nil {
			return err
		}
		if (data) == nil {
			return nil
		}

		if len(f.Field) > 0 {
			return outputFieldBatchASNDetails(data, f.Field, false, false)
		}

		if f.Yaml {
			return outputYAML(data)
		}

		return outputJSON(data)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}

	return nil
}
