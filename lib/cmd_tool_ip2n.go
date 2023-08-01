package lib

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

// CmdToolIP2nFlags are flags expected by CmdToolIP2n
type CmdToolIP2nFlags struct {
	Help    bool
	NoColor bool
}

// Init initializes the common flags available to CmdToolIP2n with sensible
func (f *CmdToolIP2nFlags) Init() {
	_h := "see description in --help"
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVar(
		&f.NoColor,
		"nocolor", false,
		_h,
	)
}

// CmdToolIP2n converts an IP address to a number
func CmdToolIP2n(f CmdToolIP2nFlags, args []string, printHelp func()) error {
	if f.NoColor {
		color.NoColor = true
	}

	if len(args) == 0 {
		printHelp()
		return nil
	}

	ipString := args[0]
	res, err := IPtoDecimalStr(ipString)
	if err != nil {
		return err
	}

	fmt.Println(res)
	return nil
}
