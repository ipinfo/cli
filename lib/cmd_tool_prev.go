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

	actionFunc := func(input string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			ActionForIPNextPrev(input, decrement)
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