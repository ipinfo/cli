package lib

import (
	"fmt"
	"net"

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
	actionFunc := func(input string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			ActionIsLinkLocalMulticast(input)
		case INPUT_TYPE_IP_RANGE:
			ActionIsLinkLocalMulticastRange(input)
		case INPUT_TYPE_CIDR:
			ActionIsLinkLocalMulticastCIDR(input)
		}
		return nil
	}
	err := GetInputFrom(args, true, true, actionFunc)
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
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		fmt.Println("Invalid Range Input", err)
		return
	}
	ipStart := net.ParseIP(ipRange.Start)

	isLinkLocalMulticast := ipStart.IsLinkLocalMulticast()

	fmt.Printf("%s,%v\n", input, isLinkLocalMulticast)
}
func ActionIsLinkLocalMulticastCIDR(input string) {
	_, ipNet, err := net.ParseCIDR(input)
	if err != nil {
		fmt.Println("Invalid CIDR input", err)
		return
	}
	isLinkLocalMulticast := ipNet.IP.IsLinkLocalMulticast()

	fmt.Printf("%s,%v\n", input, isLinkLocalMulticast)
}
