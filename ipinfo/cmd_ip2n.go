package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsIP2n = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--nocolor": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
		"-f":        predict.Set(predictReadFmts),
		"--format":  predict.Set(predictReadFmts),
	},
}

func printHelpIp2n() {

	fmt.Printf(
		`Usage: %s ip2n <ip>

Example:
  %s ip2n "190.87.89.1"
  %s ip2n "2001:0db8:85a3:0000:0000:8a2e:0370:7334
  %s ip2n "2001:0db8:85a3::8a2e:0370:7334
  %s ip2n "::7334
  %s ip2n "7334::""
	

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase, progBase, progBase, progBase, progBase, progBase)
}

func ip2nHelp() (err error) {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable colored output.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpDefault()
		return nil
	}

	// currently we do nothing by default.
	printHelpIp2n()
	return nil
}

func cmdIP2n() error {
	var err error
	var res string
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	if strings.TrimSpace(cmd) == "" {
		err := ip2nHelp()
		if err != nil {
			return err
		}
		return nil
	}

	res, err = calcIP2n(cmd)

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
		err := ip2nHelp()
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Println(res)

	return nil
}

func calcIP2n(strIP string) (string, error) {
	if isIPv6Address(strIP) {
		ip := net.ParseIP(strIP)
		if ip == nil {
			fmt.Println("Invalid IPv6 address")
			return "", errors.New("invalid IPv6 address: '" + strIP + "'")
		}

		decimalIP := IP6toSInt(ip)
		return decimalIP.String(), nil
	}
	if isIPv4Address(strIP) {
		ip := net.ParseIP(strIP)
		if ip == nil {
			return "", errors.New("invalid IPv4 address: '" + strIP + "'")
		}
		return strconv.FormatInt(IP4toInt(ip), 10), nil
	} else {
		return "", errors.New("invalid IP address: '" + strIP + "'")
	}
}
