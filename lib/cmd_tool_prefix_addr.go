package lib

import (
	"fmt"
	"net/netip"

	"github.com/spf13/pflag"
)

type CmdToolPrefixAddrFlags struct {
	Help bool
}

func (f *CmdToolPrefixAddrFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdToolPrefixAddr(f CmdToolPrefixAddrFlags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_CIDR:
			prefix, err := netip.ParsePrefix(input)
			if err != nil {
				return err
			}
			fmt.Printf("Address: %s\n", prefix.Addr())
		default:
			fmt.Printf("%s,%v\n", input, false)
		}
		return nil
	}

	return GetInputFrom(args, true, true, op)
}
