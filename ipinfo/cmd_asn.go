package main

import (
	"errors"
	"fmt"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
	"os"
	"strings"
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

func cmdASNDefault() error {
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

	var asns []string

	op := func(string string, inputType lib.INPUT_TYPE) error {
		switch inputType {
		case lib.INPUT_TYPE_ASN:
			asns = append(asns, strings.ToUpper(string))
		default:
			return lib.ErrInvalidInput
		}
		return nil
	}

	err := lib.ProcessStringsFromStdin(op)

	ii = prepareIpinfoClient(f.Token)
	if ii.Token == "" {
		return errors.New("bulk lookups require a token; login via `ipinfo init`.")
	}

	data, err := ii.GetASNDetailsBatch(asns, ipinfo.BatchReqOpts{
		TimeoutPerBatch:              60 * 30, // 30min
		ConcurrentBatchRequestsLimit: 20,
	})
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
		err = cmdASNDefault()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}

	return nil
}
