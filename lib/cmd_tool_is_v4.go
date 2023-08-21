package lib

import (
	"fmt"
	"net"

	"github.com/spf13/pflag"
)

type CmdToolIs_v4Flags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIs_v4Flags) Init() {
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

func CmdToolIs_v4(
	f CmdToolIs_v4Flags,
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
			ActionForIsV4(input)
		case INPUT_TYPE_IP_RANGE:
			ActionForIsV4Range(input)
		case INPUT_TYPE_CIDR:
			ActionForIsV4CIDR(input)
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

func ActionForIsV4(input string) {
	ip := net.ParseIP(input)
	if ip != nil && ip.To4() != nil {
		fmt.Printf("%s is an IPv4 input\n", input)
	} else {
		fmt.Printf("%s is not an IPv4 input\n", input)
	}
}

func ActionForIsV4Range(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err != nil {
		fmt.Println("Invalid IP range input:", err)
		return
	}
	startIP := net.ParseIP(ipRange.Start)
	if startIP != nil && startIP.To4() != nil {
		fmt.Printf("%s is an IPv4 input\n", input)
	} else {
		fmt.Printf("%s is not an IPv4 input\n", input)
	}
}

func ActionForIsV4CIDR(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err == nil {
		isCIDRIPv4 := ipnet.IP.To4() != nil
		if isCIDRIPv4 {
			fmt.Printf("%s is an IPv4 input\n", input)
		} else {
			fmt.Printf("%s is not an IPv4 input.\n", input)
		}
	} else {
		fmt.Println("Invalid CIDR input:", err)
	}
}
