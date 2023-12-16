package lib

import (
	"fmt"
	"net/netip"

	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/spf13/pflag"
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
	if f.Help {
		printHelp()
		return nil
	}

	op := func(input string, inputType ipUtils.INPUT_TYPE) error {
		isOneIp := false
		switch inputType {
		case ipUtils.INPUT_TYPE_CIDR:
			prefix, err := netip.ParsePrefix(input)
			if err != nil {
				return ipUtils.ErrInvalidInput
			}
			isOneIp = prefix.IsSingleIP()
		case ipUtils.INPUT_TYPE_IP:
			isOneIp = true
		case ipUtils.INPUT_TYPE_IP_RANGE:
			isOneIp = ipRangeContainsExactlyOneIP(input)
		default:
			return ipUtils.ErrInvalidInput
		}
		fmt.Printf("%s,%v\n", input, isOneIp)
		return nil
	}

	return ipUtils.GetInputFrom(args, true, true, op)
}

func ipRangeContainsExactlyOneIP(ipRangeStr string) bool {
	ipRange, err := ipUtils.IPRangeStrFromStr(ipRangeStr)
	if err != nil {
		return false
	}

	return ipRange.Start == ipRange.End
}
