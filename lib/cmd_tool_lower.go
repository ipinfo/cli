package lib

import (
	"fmt"
	"net"

	"github.com/ipinfo/cli/lib/iputil"
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

	actionFunc := func(input string, inputType iputil.INPUT_TYPE) error {
		var err error
		switch inputType {
		case iputil.INPUT_TYPE_IP:
			fmt.Println(input)
		case iputil.INPUT_TYPE_IP_RANGE:
			err = ActionForRange(input)
		case iputil.INPUT_TYPE_CIDR:
			err = ActionForCIDR(input)
		default:
			return iputil.ErrInvalidInput
		}
		return err
	}
	err := iputil.GetInputFrom(args, true, true, actionFunc)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionForRange(input string) error {
	ipRange, err := iputil.IPRangeStrFromStr(input)
	if err != nil {
		return err
	}
	fmt.Println(ipRange.Start)
	return nil
}

func ActionForCIDR(input string) error {
	_, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		return err
	}

	var lower string
	if ipnet.IP.To4() != nil {
		ipRange, _ := iputil.IPRangeStrFromCIDR(input)
		lower = ipRange.Start
	} else if ipnet.IP.To16() != nil {
		ipRange, _ := iputil.IP6RangeStrFromCIDR(input)
		lower = ipRange.Start
	}

	fmt.Println(lower)
	return nil
}
