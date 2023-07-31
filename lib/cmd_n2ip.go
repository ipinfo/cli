package lib

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

// CmdN2IPFlags are flags expected by CmdN2IP
type CmdN2IPFlags struct {
	Help    bool
	NoColor bool
	ipv6    bool
}

// Init initializes the common flags available to CmdN2IP with sensible
func (f *CmdN2IPFlags) Init() {
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
	pflag.BoolVarP(
		&f.ipv6,
		"ipv6", "6", false,
		_h,
	)
}

// CmdN2IP converts a number to an IP address
func CmdN2IP(f CmdN2IPFlags, args []string, printHelp func()) error {
	if f.NoColor {
		color.NoColor = true
	}

	if len(args) == 0 {
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

	// Evaluate the postfix expression
	postfix := InfixToPostfix(tokens)
	result, err := EvaluatePostfix(postfix)
	if err != nil {
		return err
	}

	// Convert to IP
	// Precision should be 0 i.e. number of digits after decimal
	// as ip cannot be derived from a float
	res, err := DecimalStrToIP(result.Text('f', 0), f.ipv6)
	if err != nil {
		return err
	}

	fmt.Println(res.String())
	return nil
}
