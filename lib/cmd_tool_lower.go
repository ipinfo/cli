package lib

import (
	"fmt"
	"net"
	"os"

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

	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0

	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}

	if isStdin {
		return scanrdr(os.Stdin, processIPRangeOrCIDRLower)
	}

	for _, input := range args {
		if err := processIPRangeOrCIDRLower(input); err != nil {
			return err
		}
	}
	return nil
}

func processIPRangeOrCIDRLower(input string) error {
	if ipRange, err := IPRangeStrFromStr(input); err == nil {
		// If it's an IP range, print the starting IP in the range.
		fmt.Printf("%s\n", ipRange.Start)
		return nil
	}

	if ip := net.ParseIP(input); ip != nil {
		// If it's a simple IP address, print the IP itself.
		fmt.Printf("%s\n", ip)
		return nil
	}

	if _, ipnet, err := net.ParseCIDR(input); err == nil {
		// If it's a CIDR, print the Starting IP address of the CIDR.
		fmt.Printf("%s\n", ipnet.IP)
		return nil
	}

	return fmt.Errorf("invalid input: %s", input)
}
