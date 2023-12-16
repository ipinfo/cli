package lib

import (
	"fmt"
	"net/netip"

	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/spf13/pflag"
)

// CmdToolIsV4In6Flags are flags expected by CmdToolIsV4In6
type CmdToolIsV4In6Flags struct {
	Help bool
}

// Init initializes the common flags available to CmdToolIsV4In6 with sensible
func (f *CmdToolIsV4In6Flags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdToolIsV4In6 checks if given ip is an IPv4-mapped IPv6 address.
func CmdToolIsV4In6(f CmdToolIsV4In6Flags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, inputType ipUtils.INPUT_TYPE) error {
		switch inputType {
		case ipUtils.INPUT_TYPE_IP:
			addr, err := netip.ParseAddr(input)
			if err != nil {
				return ipUtils.ErrInvalidInput
			}

			fmt.Printf("%s,%t\n", input, addr.Is4In6())
		default:
			return ipUtils.ErrInvalidInput
		}
		return nil
	}

	return ipUtils.GetInputFrom(args, true, true, op)
}
