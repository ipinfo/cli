package lib

import (
	"os"

	"github.com/spf13/pflag"
)

// CmdPripsFlags are flags expected by CmdPrips.
type CmdPripsFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdPrips with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdPripsFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdPrips is the common core logic for the prips command.
func CmdPrips(
	f CmdPripsFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	// require args and/or stdin.
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}

	return IPListWriteFromAllSrcs(args)
}
