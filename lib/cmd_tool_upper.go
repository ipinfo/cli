package lib

import (
	"fmt"
	"net"
	"os"

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

	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0

	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}

	if isStdin {
		return scanrdr(os.Stdin, processIPRangeOrCIDRUpper)
	}

	for _, input := range args {
		if err := processIPRangeOrCIDRUpper(input); err != nil {
			return err
		}
	}
	return nil
}

func processIPRangeOrCIDRUpper(input string) error {
	ipRange, err := IPRangeStrFromStr(input)
	if err == nil {
		fmt.Println(ipRange.End)
		return nil
	}

	if ip := net.ParseIP(input); ip != nil {
		fmt.Println(input)
		return nil
	}

	if _, ipnet, err := net.ParseCIDR(input); err == nil {
		ipRange, err := IPRangeStrFromCIDR(ipnet.String())
		if err == nil {
			fmt.Println(ipRange.End)
			return nil
		}
	}

	fmt.Printf("Error parsing input: %v\n", err)
	return err
}
