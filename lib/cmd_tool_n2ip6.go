package lib

import (
	"errors"
	"fmt"

	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/spf13/pflag"
)

// CmdToolN2IP6Flags are flags expected by CmdToolN2IP6
type CmdToolN2IP6Flags struct {
	Help bool
}

// Init initializes the common flags available to CmdToolN2IP6 with sensible
func (f *CmdToolN2IP6Flags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdToolN2IP6 converts a number to an IPv6 address
func CmdToolN2IP6(f CmdToolN2IP6Flags, args []string, printHelp func()) error {
	if len(args) == 0 || f.Help {
		printHelp()
		return nil
	}

	expression := args[0]
	if IsInvalidInfix(expression) {
		return ipUtils.ErrInvalidInput
	}

	// NOTE: n2ip6 also accepts an expression, hence the tokenization and evaluation.

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
	res, err := ipUtils.DecimalStrToIP(result.Text('f', 0), true)
	if err != nil {
		return errors.New("number is too large")
	}

	fmt.Println(res.String())
	return nil
}
