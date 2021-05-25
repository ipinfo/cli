package lib

import (
	"os"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// CompletionsCIDR2Range are the completions for the cidr2range command.
var CompletionsCIDR2Range = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

// CmdCIDR2RangeFlags are flags expected by CmdCIDR2Range.
type CmdCIDR2RangeFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdCIDR2Range with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdCIDR2RangeFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdCIDR2Range is the common core logic for the cidr2range command.
func CmdCIDR2Range(
	f CmdCIDR2RangeFlags,
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
