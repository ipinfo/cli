package lib

import (
	"os"

	"github.com/ipinfo/cli/lib/iputil"
	"github.com/spf13/pflag"
)

// CmdCIDR2IPFlags are flags expected by CmdCIDR2IP.
type CmdCIDR2IPFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdCIDR2IP with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdCIDR2IPFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdCIDR2IP(f CmdCIDR2IPFlags, args []string, printHelp func()) error {
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

	return iputil.IPListWriteFrom(args, true, false, false, true, true)
}
