package lib

import (
	"fmt"

	"github.com/ipinfo/cli/lib/iputil"
	"github.com/spf13/pflag"
)

// CmdToolIP2nFlags are flags expected by CmdToolIP2n
type CmdToolIP2nFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdToolIP2n with sensible
func (f *CmdToolIP2nFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdToolIP2n converts an IP address to a number
func CmdToolIP2n(f CmdToolIP2nFlags, args []string, printHelp func()) error {
	if len(args) == 0 || f.Help {
		printHelp()
		return nil
	}

	ipString := args[0]
	res, err := iputil.IPtoDecimalStr(ipString)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}
