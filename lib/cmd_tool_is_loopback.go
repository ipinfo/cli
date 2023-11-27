package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolIsLoopbackFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIsLoopbackFlags) Init() {
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

func CmdToolIsLoopback(
	f CmdToolIsLoopbackFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	actionFuncLoopBack := func(input string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			ActionForIsLoopBack(input)
		case INPUT_TYPE_IP_RANGE:
			ActionForISLoopBackRange(input)
		case INPUT_TYPE_CIDR:
			ActionForIsLoopBackCIDR(input)
		}
		return nil
	}
	err := GetInputFrom(args, true, true, actionFuncLoopBack)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionForIsLoopBack(input string) {
	ip := net.ParseIP(input)
	isLoopBack := ip.IsLoopback()

	fmt.Printf("%s,%v\n", input, isLoopBack)
}

func ActionForISLoopBackRange(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		return
	}

	ipStart := net.ParseIP(ipRange.Start)
	isLoopBack := ipStart.IsLoopback()

	fmt.Printf("%s,%v\n", input, isLoopBack)
}

func ActionForIsLoopBackCIDR(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		return
	}

	isCIDRLoopBack := ipnet.IP.IsLoopback()

	fmt.Printf("%s,%v\n", input, isCIDRLoopBack)
}
