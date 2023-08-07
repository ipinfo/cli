package lib

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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
	stdin bool,
	ip bool,
	iprange bool,
	cidr bool,
) error {
	if f.Help {
		printHelp()
		return nil
	}

	actionStdin := func(input string, iprange, cidr bool) {
		ActionForStdinUpper(input, iprange, cidr)
	}
	actionRange := func(input string) {
		ActionForRangeUpper(input)
	}
	actionCidr := func(input string) {
		ActionForCIDRUpper(input)
	}
	actionFile := func(input string, iprange, cidr bool) {
		ActionForFileUpper(input, iprange, cidr)
	}

	// Process inputs using the IPInputAction function.
	err := IPInputAction(args, true, true, true, true, true,
		actionStdin, actionRange, actionCidr, actionFile)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionForStdinUpper(input string, iprange bool, cidr bool) {
	ip := net.ParseIP(input)
	if ip != nil {
		fmt.Println(ip)
	} else if iprange {
		ActionForRangeUpper(input)
	} else if cidr {
		ActionForCIDRUpper(input)
	}
}

func ActionForRangeUpper(input string) {
	ipRange, err := IPRangeStrFromStr(input)
	if err == nil {
		fmt.Println(ipRange.End)
	}
}

func ActionForCIDRUpper(input string) {
	_, ipnet, err := net.ParseCIDR(input)
	if err == nil {
		ipRange, err := IPRangeStrFromCIDR(ipnet.String())
		if err == nil {
			fmt.Println(ipRange.End)
		}
	}
}

func ActionForFileUpper(pathToFile string, iprange bool, cidr bool) {
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
		ActionForStdinUpper(input, iprange, cidr)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}
