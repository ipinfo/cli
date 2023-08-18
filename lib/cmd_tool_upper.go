package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolUpperFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolUpperFlags) Init() {
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

func CmdToolUpper(
	f CmdToolUpperFlags,
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
			ActionForIPUpper(input)
		case INPUT_TYPE_IP_RANGE:
			ActionForRangeUpper(input)
		case INPUT_TYPE_CIDR:
			ActionForCIDRUpper(input)
		default:
			return ErrNotIP
		}
		return nil
	}
    err := GetInputFrom(args, true, true, actionFunc)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func ActionForIPUpper(input string) {
	ip := net.ParseIP(input)
	if ip != nil {
		fmt.Println(ip)
	}
}

func ActionForRangeUpper(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ipRange.End) // Print the end IP of the range
}

func ActionForCIDRUpper(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	ipRange, err := IPRangeStrFromCIDR(ipnet.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ipRange.End) // Print the end IP of the CIDR
}
