package lib

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

// CmdN2IP6Flags are flags expected by CmdN2IP6
type CmdN2IP6Flags struct {
	Help    bool
	NoColor bool
}

// Init initializes the common flags available to CmdN2IP6 with sensible
func (f *CmdN2IP6Flags) Init() {
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

// CmdN2IP6 converts a number to an IPv6 address
func CmdN2IP6(f CmdN2IP6Flags, args []string, printHelp func()) error {
	if f.NoColor {
		color.NoColor = true
	}

	if f.Help {
		printHelp()
		return nil
	}

	expression := args[0]

	if IsInvalidInfix(expression) {
		return ErrInvalidInput
	}

	// n2ip also accepts an expression which is why the following
	// Steps are being done
	// Convert to postfix
	// If it is a single number and not an expression
	// The tokenization and evaluation would have no effect on the number

	// Tokenize the expression
	tokens, err := TokenizeInfix(expression)
	if err != nil {
		return err
	}

	postfix := InfixToPostfix(tokens)
	//
	// Evaluate the postfix expression
	result, err := EvaluatePostfix(postfix)
	if err != nil {
		return err
	}

	// Convert to IP
	// Precision should be 0 i.e. number of digits after decimal
	// as ip cannot be derived from a float
	res, err := DecimalStrToIP(result.Text('f', 0), true)
	if err != nil {
		return err
	}

	fmt.Println(res.String())
	return nil
}
