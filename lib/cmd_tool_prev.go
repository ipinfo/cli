package lib

import (
	"fmt"

	"github.com/ipinfo/cli/lib/iputil"
	"github.com/spf13/pflag"
)

type CmdToolPrevFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolPrevFlags) Init() {
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

func CmdToolPrev(
	f CmdToolPrevFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	actionFunc := func(input string, inputType iputil.INPUT_TYPE) error {
		switch inputType {
		case iputil.INPUT_TYPE_IP:
			newIP := iputil.IpAdd(input, -1)
			fmt.Println(newIP)
		}
		return nil
	}
	err := iputil.GetInputFrom(args, true, true, actionFunc)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
