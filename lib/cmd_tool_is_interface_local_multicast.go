package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolIsInterfaceLocalMulticastFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIsInterfaceLocalMulticastFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help",
	)
	pflag.BoolVarP(
		&f.Quiet,
		"quiet", "q", false,
		"quiet mode; suppress additional output.",
	)
}

func CmdToolIsInterfaceLocalMulticast(
	f CmdToolIsInterfaceLocalMulticastFlags,
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
			ActionIsInterfaceLocalMulticast(input)
		case INPUT_TYPE_IP_RANGE:
			ActionIsInterfaceLocalMulticastRange(input)
		case INPUT_TYPE_CIDR:
			ActionIsInterfaceLocalMulticastCIDR(input)
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

func ActionIsInterfaceLocalMulticast(input string) {
	ip := net.ParseIP(input)
	isInterfaceLocalMulticast := ip.IsInterfaceLocalMulticast()

	fmt.Printf("%s,%v\n", input, isInterfaceLocalMulticast)
}

func ActionIsInterfaceLocalMulticastRange(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		return
	}

	ipStart := net.ParseIP(ipRange.Start)
	isInterfaceLocalMulticast := ipStart.IsInterfaceLocalMulticast()

	fmt.Printf("%s,%v\n", input, isInterfaceLocalMulticast)
}

func ActionIsInterfaceLocalMulticastCIDR(input string) {
	_, ipNet, err := net.ParseCIDR(input)
	if err != nil {
		return
	}

	isInterfaceLocalMulticast := ipNet.IP.IsInterfaceLocalMulticast()

	fmt.Printf("%s,%v\n", input, isInterfaceLocalMulticast)
}
