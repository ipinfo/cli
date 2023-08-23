package lib

import (
	"fmt"
	"github.com/spf13/pflag"
)

// CmdToolIsV6Flags are flags expected by CmdToolIsV6.
type CmdToolIsV6Flags struct {
	Help bool
}

// Init initializes the common flags available to CmdToolIsV6 with sensible
func (f *CmdToolIsV6Flags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdToolIsV6(f CmdToolIsV6Flags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, input_type INPUT_TYPE) error {
		fmt.Printf("%s %t\n", input, StrIsIPv6Str(input))
		return nil
	}

	return GetInputFrom(args, true, true, op)
}
