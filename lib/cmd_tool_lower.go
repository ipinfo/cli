package lib

import (
	"fmt"
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

	processCIDR := func(cidrStr string) error {
		ipRange, err := IPRangeStrFromCIDR(cidrStr)
		if err != nil {
			if !f.Quiet {
				fmt.Printf("Error parsing CIDR: %v\n", err)
			}
			return nil
		}

		fmt.Println(ipRange.Start)
		return nil
	}

	if isStdin {
		return scanrdr(os.Stdin, processCIDR)
	}

	for _, cidrStr := range args {
		if err := processCIDR(cidrStr); err != nil {
			return err
		}
	}
	return nil
}
