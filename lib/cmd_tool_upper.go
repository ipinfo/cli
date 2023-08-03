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
	stdin bool,
	ip bool,
	iprange bool,
	cidr bool,
) error {
	if f.Help {
		printHelp()
		return nil
	}

	if stdin {
		stat, _ := os.Stdin.Stat()

		isPiped := (stat.Mode() & os.ModeNamedPipe) != 0
		isTyping := (stat.Mode()&os.ModeCharDevice) != 0 && len(args) == 0

		if isTyping {
			fmt.Println("** manual input mode **")
			fmt.Println("Enter all IPs, one per line:")
		}

		if isPiped || isTyping || stat.Size() > 0 {
			return scanrdr(os.Stdin, processIPRangeOrCIDRUpper)
		}
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
	return ErrInvalidInput
}
