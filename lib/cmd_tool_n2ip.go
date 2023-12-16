package lib

import (
	"errors"
	"fmt"

	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/spf13/pflag"
)

// CmdToolN2IPFlags are flags expected by CmdToolN2IP
type CmdToolN2IPFlags struct {
	Help bool
	ipv6 bool
}

// Init initializes the common flags available to CmdToolN2IP with sensible
func (f *CmdToolN2IPFlags) Init() {
	_h := "see description in --help"
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.ipv6,
		"ipv6", "6", false,
		_h,
	)
}

// CmdToolN2IP converts a number to an IP address
func CmdToolN2IP(f CmdToolN2IPFlags, args []string, printHelp func()) error {
	if len(args) == 0 || f.Help {
		printHelp()
		return nil
	}

	expression := args[0]
	if IsInvalidInfix(expression) {
		return ipUtils.ErrInvalidInput
	}

	// NOTE: n2ip also accepts an expression, hence the tokenization and evaluation.

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
	res, err := ipUtils.DecimalStrToIP(result.Text('f', 0), f.ipv6)
	if err != nil {
		return errors.New("number is too large")
	}

	fmt.Println(res.String())
	return nil
}
