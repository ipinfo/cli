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
	stdin bool,
	ip bool,
	iprange bool,
	cidr bool,
) error {
	if f.Help {
		printHelp()
		return nil
	}

	if !stdin && !ip && !iprange && !cidr {
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
			return scanrdr(os.Stdin, processIPRangeOrCIDRLower)
		}
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
		fmt.Println(ipRange.Start)
		return nil
	}

	if ip := net.ParseIP(input); ip != nil {
		// If it's a simple IP address, print the IP itself
		fmt.Printf("%s\n", ip)
		return nil
	}

	if _, network, err := net.ParseCIDR(input); err == nil {
		// If it's a CIDR, print the Starting IP address of the CIDR.
		fmt.Println(network.IP)
		return nil
	}
	return ErrInvalidInput
}
