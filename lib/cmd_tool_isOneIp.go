package lib

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
)

// CmdToolIsOneIpFlags are flags expected by CmdToolIP2n
type CmdToolIsOneIpFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdToolIsOneIp with sensible
func (f *CmdToolIsOneIpFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

func CmdToolIsOneIp(f CmdToolIsOneIpFlags, args []string, printHelp func()) error {
	if len(args) == 0 || f.Help {
		printHelp()
		return nil
	}

	op := func(input string, inputType INPUT_TYPE) error {
		isOneIp := false
		switch inputType {
		case INPUT_TYPE_CIDR:
			isOneIp = CIDRContainsExactlyOneIP(input)
		case INPUT_TYPE_IP:
			isOneIp = true
		case INPUT_TYPE_IP_RANGE:
			isOneIp = ipRangeContainsExactlyOneIP(input)
		default:
			return ErrInvalidInput
		}
		fmt.Printf("%s %v\n", input, isOneIp)
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

func ipRangeContainsExactlyOneIP(ipRangeStr string) bool {
	ipRange, err := IPRangeStrFromStr(ipRangeStr)
	if err != nil {
		return false
	}

	return ipRange.Start == ipRange.End
}
