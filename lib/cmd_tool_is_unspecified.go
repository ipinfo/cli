package lib

import (
	"fmt"
	"net"

	"github.com/ipinfo/cli/lib/iputil"
	"github.com/spf13/pflag"
)

type CmdToolIsUnspecifiedFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolIsUnspecifiedFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.Quiet,
		"quiet", "q", false,
		"quiet mode;suppress additional output.",
	)
}

func CmdToolIsUnspecified(
	f CmdToolIsUnspecifiedFlags,
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
			ActionIsUnspecified(input)
		}
		return nil
	}
	err := iputil.GetInputFrom(args, true, true, actionFunc)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ActionIsUnspecified(input string) {
	ip := net.ParseIP(input)
	isUnspecified := ip.IsUnspecified()

	fmt.Printf("%s,%v\n", input, isUnspecified)
}
