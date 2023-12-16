package lib

import (
	"fmt"
	"net"

	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/spf13/pflag"
)

type CmdToolIsLinkLocalUnicastFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIsLinkLocalUnicastFlags) Init() {
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

func CmdToolIsLinkLocalUnicast(
	f CmdToolIsLinkLocalUnicastFlags,
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
			ActionIsLinkLocalUnicast(input)
		case ipUtils.INPUT_TYPE_IP_RANGE:
			ActionIsLinkLocalUnicastRange(input)
		case ipUtils.INPUT_TYPE_CIDR:
			ActionIsLinkLocalUnicastCIDR(input)
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

func ActionIsLinkLocalUnicast(input string) {
	ip := net.ParseIP(input)
	isLinkLocalUnicast := ip.IsLinkLocalUnicast()

	fmt.Printf("%s,%v\n", input, isLinkLocalUnicast)
}

func ActionIsLinkLocalUnicastRange(input string) {
	ipRange, err := ipUtils.IPRangeStrFromStr(input)
	if err != nil {
		return
	}

	ipStart := net.ParseIP(ipRange.Start)
	isLinkLocalUnicast := ipStart.IsLinkLocalUnicast()

	fmt.Printf("%s,%v\n", input, isLinkLocalUnicast)
}

func ActionIsLinkLocalUnicastCIDR(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		return
	}

	isLinkLocalUnicast := ipnet.IP.IsLinkLocalUnicast()

	fmt.Printf("%s,%v\n", input, isLinkLocalUnicast)
}
