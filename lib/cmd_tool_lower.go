package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolLowerFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolLowerFlags) Init() {
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

func CmdToolLower(
	f CmdToolLowerFlags,
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
			ActionForIP(input)
		case INPUT_TYPE_IP_RANGE:
			ActionForRange(input)
		case INPUT_TYPE_CIDR:
			ActionForCIDR(input)
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

func ActionForIP(input string) {
	fmt.Println(input)
}

func ActionForRange(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ipRange.Start)
}

func ActionForCIDR(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ipnet.IP)
}
