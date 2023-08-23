package lib

import (
	"fmt"
	"github.com/spf13/pflag"
)

// CmdToolIsV4Flags are flags expected by CmdToolIsV4.
type CmdToolIsV4Flags struct {
	Help bool
}

// Init initializes the common flags available to CmdToolIsV4 with sensible
func (f *CmdToolIsV4Flags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdToolIsV4(f CmdToolIsV4Flags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, input_type INPUT_TYPE) error {
		fmt.Printf("%s %t\n", input, StrIsIPv4Str(input))
		return nil
	}

	return GetInputFrom(args, true, true, op)
}
