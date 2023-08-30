package lib

import (
	"fmt"
	"github.com/spf13/pflag"
)

// CmdToolUnmapFlags are flags expected by CmdToolUnmap
type CmdToolUnmapFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdToolUnmap with sensible
func (f *CmdToolUnmapFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdToolUnmap converts a number to an IP address
func CmdToolUnmap(f CmdToolUnmapFlags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, input_type INPUT_TYPE) error {
		switch input_type {
		case INPUT_TYPE_IP:
			ip, err := IPtoDecimalStr(input)
			if err != nil {
				return err
			}
			fmt.Println(ip)
		default:
			return ErrNotIP
		}
		return nil
	}

	return GetInputFrom(args, true, true, op)
}
