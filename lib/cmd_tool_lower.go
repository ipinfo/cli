package lib

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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

	actionStdin := func(input string, iprange, cidr bool) {
		ActionForStdinLower(input, iprange, cidr)
	}
	actionRange := func(input string) {
		ActionForRangeLower(input)
	}
	actionCidr := func(input string) {
		ActionForCIDRLower(input)
	}
	actionFile := func(input string, iprange, cidr bool) {
		ActionForFileLower(input, iprange, cidr)
	}

	// Process inputs using the IPInputAction function.
	err := IPInputAction(args, true, true, true, true, true,
		actionStdin, actionRange, actionCidr, actionFile)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionForStdinLower(input string, iprange bool, cidr bool) {
	ip := net.ParseIP(input)
	if ip != nil {
		fmt.Println(ip)
	} else if iprange {
		ActionForRangeLower(input)
	} else if cidr {
		ActionForCIDRLower(input)
	}
}

func ActionForRangeLower(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err == nil {
		fmt.Println(ipRange.Start)
	}
}

func ActionForCIDRLower(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err == nil {
		fmt.Println(ipnet.IP)
	}
}

func ActionForFileLower(pathToFile string, iprange bool, cidr bool) {
	f, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		ActionForStdinLower(input, iprange, cidr) 
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
