package lib

import (
	"fmt"
	"net/netip"

	"github.com/ipinfo/cli/lib/ipUtils"
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

	op := func(input string, input_type ipUtils.INPUT_TYPE) error {
		switch input_type {
		case ipUtils.INPUT_TYPE_IP:
			addr, err := netip.ParseAddr(input)
			if err != nil {
				return err
			}
			fmt.Println(addr.Unmap())
		}
		return nil
	}

	return ipUtils.GetInputFrom(args, true, true, op)
}
