package lib

import (
	"os"

	"github.com/spf13/pflag"
)

// CmdRange2IPFlags are flags expected by CmdRange2IP.
type CmdRange2IPFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdRange2IP with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdRange2IPFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdRange2IP(f CmdRange2IPFlags, args []string, printHelp func()) error {
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

	return IPListWriteFrom(args, true, true, true, false, true)
}
