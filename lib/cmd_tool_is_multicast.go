package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolIsMulticastFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIsMulticastFlags) Init() {
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

func CmdToolIsMulticast(
	f CmdToolIsMulticastFlags,
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
			ActionForIsMulticast(input)
		case INPUT_TYPE_IP_RANGE:
			ActionForIsMulticastRange(input)
		case INPUT_TYPE_CIDR:
			ActionForIsMulticastCIDR(input)
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

func ActionForIsMulticast(input string) {
	ip := net.ParseIP(input)
	isMulticast := ip.IsMulticast()

	fmt.Printf("%s,%v\n", input, isMulticast)
}

func ActionForIsMulticastRange(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		return
	}

	ipStart := net.ParseIP(ipRange.Start)
	isMulticast := ipStart.IsMulticast()

	fmt.Printf("%s,%v\n", ipStart, isMulticast)
}

func ActionForIsMulticastCIDR(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		return
	}

	isCIDRMulticast := ipnet.IP.IsMulticast()

	fmt.Printf("%s,%v\n", input, isCIDRMulticast)
}
