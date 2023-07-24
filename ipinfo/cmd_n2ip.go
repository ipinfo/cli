package main

import (
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var forceIpv6 bool

var completionsN2IP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-f":        predict.Set(predictReadFmts),
		"--format":  predict.Set(predictReadFmts),
		"-6":        predict.Set(predictReadFmts),
		"--ipv6":    predict.Set(predictReadFmts),
	},
}

func printHelpN2IP() {

	fmt.Printf(
		`Usage: %s n2ip [<opts>] <expr>

Example:
  %s n2ip "2*2828-1"
  %s n2ip "190.87.89.1*2"
  %s n2ip "2001:0db8:85a3:0000:0000:8a2e:0370:7334*6"

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
    --ipv6, -6
      force conversion to IPv6 address
`, progBase, progBase, progBase, progBase)
}

func n2ipHelp() (err error) {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable colored output.")

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpDefault()
		return nil
	}

	// currently we do nothing by default.
	printHelpN2IP()
	return nil
}

func cmdN2IP() error {
	pflag.BoolVarP(&forceIpv6, "ipv6", "6", false, "force conversion to IPv6 address")
	pflag.Parse()
	var err error

	cmd := ""

	fmt.Println(os.Args)

	// Reading input from the command line
	if forceIpv6 && len(os.Args) > 3 {
		cmd = os.Args[3]
	} else if !forceIpv6 && len(os.Args) > 2 {
		cmd = os.Args[2]
	} else {
		err := n2ipHelp()
		if err != nil {
			return err
		}
		return nil
	}

	// Validate the input
	if strings.TrimSpace(cmd) == "" {
		err := n2ipHelp()
		if err != nil {
			return err
		}
		return nil
	}

	if isInvalid(cmd) {
		return errors.New("invalid expression")
	}

	// Tokenize
	tokens, err := tokeinzeExp(cmd)

	if err != nil {
		return err
	}

	// Convert to postfix
	// If it is a single number and not an expression
	// The tokenization ad evaluation would have no effect on the number
	postfix := infixToPostfix(tokens)

	// Evaluate the postfix expression
	result, err := evaluatePostfix(postfix)

	if err != nil {
		return err
	}

	// Convert to IP
	res := decimalToIP(result.String(), forceIpv6)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		err := n2ipHelp()
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Println(res)

	return nil
}

func decimalToIP(decimal string, forceIPv6 bool) net.IP {
	// Create a new big.Int with a value of 'decimal'
	num := new(big.Int)
	num, success := num.SetString(decimal, 10)
	if !success {
		fmt.Fprintf(os.Stderr, "Error parsing the decimal string: %v\n", success)
		return nil
	}

	// Convert to IPv4 if not forcing IPv6 and 'num' is within the IPv4 range
	if !forceIPv6 && num.Cmp(big.NewInt(4294967295)) <= 0 {
		ip := make(net.IP, 4)
		b := num.Bytes()
		copy(ip[4-len(b):], b)
		return ip
	}

	// Convert to IPv6
	ip := make(net.IP, 16)
	b := num.Bytes()
	copy(ip[16-len(b):], b)
	return ip
}
