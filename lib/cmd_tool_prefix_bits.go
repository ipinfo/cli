package lib

import (
	"fmt"
	"net/netip"

	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/spf13/pflag"
)

type CmdToolPrefixBitsFlags struct {
	Help bool
}

func (f *CmdToolPrefixBitsFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdToolPrefixBits(f CmdToolPrefixBitsFlags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, inputType ipUtils.INPUT_TYPE) error {
		switch inputType {
		case ipUtils.INPUT_TYPE_CIDR:
			prefix, err := netip.ParsePrefix(input)
			if err != nil {
				return err
			}
			fmt.Printf("%s,%d\n", input, prefix.Bits())
		}
		return nil
	}

	return ipUtils.GetInputFrom(args, true, true, op)
}
