package lib

import (
	"fmt"
	"net"

	"github.com/ipinfo/cli/lib/ipUtils"
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

	actionFunc := func(input string, inputType ipUtils.INPUT_TYPE) error {
		var err error
		switch inputType {
		case ipUtils.INPUT_TYPE_IP:
			fmt.Println(input)
		case ipUtils.INPUT_TYPE_IP_RANGE:
			err = ActionForRangeUpper(input)
		case ipUtils.INPUT_TYPE_CIDR:
			err = ActionForCIDRUpper(input)
		default:
			return ipUtils.ErrInvalidInput
		}
		return err
	}
	err := ipUtils.GetInputFrom(args, true, true, actionFunc)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionForRangeUpper(input string) error {
	ipRange, err := ipUtils.IPRangeStrFromStr(input)
	if err != nil {
		return err
	}
	fmt.Println(ipRange.End)
	return nil
}

func ActionForCIDRUpper(input string) error {
	_, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		return err
	}

	var upper string
	if ipnet.IP.To4() != nil {
		ipRange, _ := ipUtils.IPRangeStrFromCIDR(input)
		upper = ipRange.End
	} else if ipnet.IP.To16() != nil {
		ipRange, _ := ipUtils.IP6RangeStrFromCIDR(input)
		upper = ipRange.End
	}

	fmt.Println(upper)
	return nil
}
