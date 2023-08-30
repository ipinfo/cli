package lib

import (
	"fmt"
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

	op := func(input string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			fmt.Printf("%s,%t\n", input, isV4in6(input))
		default:
			return ErrInvalidInput
		}
		return nil
	}

	return GetInputFrom(args, true, true, op)
}

func isV4in6(ip string) bool {
	decimalFormat, err := IPtoDecimalStr(ip)
	if err != nil {
		return false
	}

	res, err := DecimalStrToIP(decimalFormat, false)
	if err != nil {
		return false
	}

	return StrIsIPv4Str(res.String())
}
