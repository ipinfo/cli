package lib

import (
	"fmt"

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

	decrement := -1
	actionStdin := func(input string) {
		ActionForStdinNextPrev(input, decrement)
	}
	actionFile := func(input string) {
		ActionForFileNextPrev(input, decrement)
	}

	err := IPInputAction(args, true, true, false, false, true,
		actionStdin, nil, nil, actionFile)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
