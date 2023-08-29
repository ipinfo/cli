package lib

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
)

// CmdToolIsSingleIpFlags are flags expected by CmdToolIP2n
type CmdToolIsSingleIpFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdToolIsSingleIp with sensible
func (f *CmdToolIsSingleIpFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdToolIsSingleIp(f CmdToolIsSingleIpFlags, args []string, printHelp func()) error {
	if len(args) == 0 || f.Help {
		printHelp()
		return nil
	}

	op := func(input string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_CIDR:
			fmt.Printf("%s %v\n", input, CIDRContainsExactlyOneIP(input))
		default:
			return ErrNotCIDR
		}
		return nil
	}

	return GetInputFrom(args, true, true, op)
}

// CIDRContainsExactlyOneIP checks whether a CIDR contains exactly one IP
func CIDRContainsExactlyOneIP(cidrStr string) bool {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false
	}

	if ipnet.IP.To4() != nil {
		ipRange, _ := IPRangeStrFromCIDR(cidrStr)
		return ipRange.Start == ipRange.End
	} else if ipnet.IP.To16() != nil {
		ipRange, _ := IP6RangeStrFromCIDR(cidrStr)
		return ipRange.Start == ipRange.End
	}

	return false
}
