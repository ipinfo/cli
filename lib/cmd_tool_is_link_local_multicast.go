package lib

import (
	"fmt"
	"net"

	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/spf13/pflag"
)

type CmdToolIsLinkLocalMulticastFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIsLinkLocalMulticastFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.Quiet,
		"quiet", "q", false,
		"quiet mode; suppress additional output.",
	)
}

func CmdToolIsLinkLocalMulticast(
	f CmdToolIsLinkLocalMulticastFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	actionFunc := func(input string, inputType ipUtils.INPUT_TYPE) error {
		switch inputType {
		case ipUtils.INPUT_TYPE_IP:
			ActionIsLinkLocalMulticast(input)
		case ipUtils.INPUT_TYPE_IP_RANGE:
			ActionIsLinkLocalMulticastRange(input)
		case ipUtils.INPUT_TYPE_CIDR:
			ActionIsLinkLocalMulticastCIDR(input)
		}
		return nil
	}
	err := ipUtils.GetInputFrom(args, true, true, actionFunc)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ActionIsLinkLocalMulticast(input string) {
	ip := net.ParseIP(input)
	isLinkLocalMulticast := ip.IsLinkLocalMulticast()

	fmt.Printf("%s,%v\n", input, isLinkLocalMulticast)
}

func ActionIsLinkLocalMulticastRange(input string) {
	ipRange, err := ipUtils.IPRangeStrFromStr(input)
	if err != nil {
		return
	}

	ipStart := net.ParseIP(ipRange.Start)
	isLinkLocalMulticast := ipStart.IsLinkLocalMulticast()

	fmt.Printf("%s,%v\n", input, isLinkLocalMulticast)
}

func ActionIsLinkLocalMulticastCIDR(input string) {
	_, ipNet, err := net.ParseCIDR(input)
	if err != nil {
		return
	}

	isLinkLocalMulticast := ipNet.IP.IsLinkLocalMulticast()

	fmt.Printf("%s,%v\n", input, isLinkLocalMulticast)
}
