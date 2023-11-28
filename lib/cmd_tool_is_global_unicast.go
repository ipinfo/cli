package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolIsGlobalUnicastFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIsGlobalUnicastFlags) Init() {
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

func CmdToolIsGlobalUnicast(
	f CmdToolIsGlobalUnicastFlags,
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
			ActionIsGlobalUnicast(input)
		case INPUT_TYPE_IP_RANGE:
			ActionIsGlobalUnicastRange(input)
		case INPUT_TYPE_CIDR:
			ActionIsGlobalUnicastCIDR(input)
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

func ActionIsGlobalUnicast(input string) {
	ip := net.ParseIP(input)
	isGlobalUnicast := ip.IsGlobalUnicast()

	fmt.Printf("%s,%v\n", input, isGlobalUnicast)
}

func ActionIsGlobalUnicastRange(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		return
	}

	ipStart := net.ParseIP(ipRange.Start)
	isGlobalUnicast := ipStart.IsGlobalUnicast()

	fmt.Printf("%s,%v\n", input, isGlobalUnicast)
}

func ActionIsGlobalUnicastCIDR(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err != nil {
		return
	}

	isGlobalUnicast := ipnet.IP.IsGlobalUnicast()

	fmt.Printf("%s,%v\n", input, isGlobalUnicast)
}
