package lib

import (
	"os"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// CompletionsRange2CIDR are the completions for the range2cidr command.
var CompletionsRange2CIDR = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// CmdRange2CIDRFlags are flags expected by CmdRange2CIDR.
type CmdRange2CIDRFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdRange2CIDR with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdRange2CIDRFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdRange2CIDR is the common core logic for the range2cidr command.
func CmdRange2CIDR(
	f CmdRange2CIDRFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	// require args.
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}

	// TODO impl logic

	return nil
}
