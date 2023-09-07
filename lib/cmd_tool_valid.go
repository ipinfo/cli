package lib

import (
	"fmt"
	"github.com/spf13/pflag"
)

// CmdToolIsValidFlags are flags expected by CmdToolIsValid
type CmdToolIsValidFlags struct {
	Help bool
	ipv6 bool
}

// Init initializes the common flags available to CmdToolIsValid with sensible
func (f *CmdToolIsValidFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdToolIsValid converts a number to an IP address
func CmdToolIsValid(f CmdToolIsValidFlags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, input_type INPUT_TYPE) error {
		switch input_type {
		case INPUT_TYPE_IP:
			fmt.Printf("%s,%v\n", input, true)
		default:
			fmt.Printf("%s,%v\n", input, false)
		}
		return nil
	}

	return GetInputFrom(args, true, true, op)
}
