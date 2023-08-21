package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolIs_v6Flags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIs_v6Flags) Init() {
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

func CmdToolIs_v6(
	f CmdToolIs_v6Flags,
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
			ActionForIsV6(input)
		case INPUT_TYPE_IP_RANGE:
			ActionForIsV6Range(input)
		case INPUT_TYPE_CIDR:
			ActionForIsV6CIDR(input)
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

func ActionForIsV6(input string) {
	ip := net.ParseIP(input)
	isIPv6 := ip != nil && ip.To16() != nil && ip.To4() == nil

	if isIPv6 {
		fmt.Printf("%s is an IPv6 input.\n", input)
	} else {
		fmt.Printf("%s is not an IPv6 input.\n", input)
	}
}

func ActionForIsV6Range(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		fmt.Println("Invalid IP range input:", err)
		return
	}

	startIP := net.ParseIP(ipRange.Start)
	isIPv6 := startIP != nil && startIP.To16() != nil && startIP.To4() == nil

	if isIPv6 {
		fmt.Printf("%s is an IPv6 input.\n", input)
	} else {
		fmt.Printf("%s is not an IPv6 input.\n", input)
	}
}

func ActionForIsV6CIDR(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err == nil {
		isCIDRIPv6 := ipnet.IP.To16() != nil && ipnet.IP.To4() == nil
		if isCIDRIPv6 {
			fmt.Printf("%s is an IPv6 input.\n", input)
		} else {
			fmt.Printf("%s is not an IPv6 input.\n", input)
		}
	} else {
		fmt.Println("Invalid CIDR input:", err)
	}
}
