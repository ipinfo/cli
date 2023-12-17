package lib

import (
	"fmt"
	"net/netip"

	"github.com/ipinfo/cli/lib/iputil"
	"github.com/spf13/pflag"
)

type CmdToolPrefixIsValidFlags struct {
	Help bool
}

func (f *CmdToolPrefixIsValidFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdToolPrefixIsValid(f CmdToolPrefixIsValidFlags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, inputType iputil.INPUT_TYPE) error {
		switch inputType {
		case iputil.INPUT_TYPE_CIDR:
			prefix, err := netip.ParsePrefix(input)
			if err != nil {
				return err
			}
			fmt.Printf("%s,%t\n", input, prefix.IsValid())
		}
		return nil
	}

	return iputil.GetInputFrom(args, true, true, op)
}
