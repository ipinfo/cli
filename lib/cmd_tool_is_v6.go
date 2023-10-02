package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolIsV6Flags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIsV6Flags) Init() {
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

func CmdToolIsV6(
	f CmdToolIsV6Flags,
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
	isIPv6 := IsIPv6(ip)

	fmt.Printf("%s,%v\n", input, isIPv6)
}

func ActionForIsV6Range(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		fmt.Println("Invalid IP range input:", err)
		return
	}

	startIP := net.ParseIP(ipRange.Start)
	isIPv6 := IsIPv6(startIP)

	fmt.Printf("%s,%v\n", input, isIPv6)
}

func ActionForIsV6CIDR(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err == nil {
		isCIDRIPv6 := IsCIDRIPv6(ipnet)
		fmt.Printf("%s,%v\n", input, isCIDRIPv6)
	} else {
		fmt.Println("Invalid CIDR input:", err)
	}
}

func IsIPv6(ip net.IP) bool {
	return ip != nil && ip.To16() != nil && ip.To4() == nil
}

func IsCIDRIPv6(ipnet *net.IPNet) bool {
	return ipnet != nil && ipnet.IP.To16() != nil && ipnet.IP.To4() == nil
}
