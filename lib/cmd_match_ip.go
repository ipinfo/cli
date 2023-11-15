package lib

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type CmdMatchIPFlags struct {
	Expression []string
	Help       bool
}

func (f *CmdMatchIPFlags) Init() {
	pflag.StringSliceVarP(
		&f.Expression,
		"expression", "e", nil,
		"IPs, subnets to be filtered.",
	)
	pflag.BoolVarP(
		&f.Help,
		"help", "", false,
		"show help.",
	)
}

func CmdMatchIP(
	f CmdMatchIPFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help || len(f.Expression) == 0 || args[0] == "" {
		printHelp()
		return nil
	}

	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0

	var sourceCIDRs, filterCIDRs []SubnetPair
	var sourceIPs, filterIPs []net.IP

	parseCIDRs := func(cidrs []string) []net.IPNet {
		parsedCIDRs := make([]net.IPNet, 0)
		for _, cidrStr := range cidrs {
			_, ipNet, err := net.ParseCIDR(cidrStr)
			if err != nil {
				continue
			}
			parsedCIDRs = append(parsedCIDRs, *ipNet)
		}

		return parsedCIDRs
	}

	source_op := func(s string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			sourceIPs = append(sourceIPs, net.ParseIP(s))
		case INPUT_TYPE_CIDR:
			pair := SubnetPair{
				Raw:    s,
				Parsed: parseCIDRs([]string{s}),
			}
			sourceCIDRs = append(sourceCIDRs, pair)
		case INPUT_TYPE_IP_RANGE:
			cidrs, err := rangeToCidrs(s)
			if err == nil {
				pair := SubnetPair{
					Raw:    s,
					Parsed: parseCIDRs(cidrs),
				}
				sourceCIDRs = append(sourceCIDRs, pair)
			}
		default:
			return ErrInvalidInput
		}
		return nil
	}

	filter_op := func(s string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			filterIPs = append(filterIPs, net.ParseIP(s))
		case INPUT_TYPE_CIDR:
			pair := SubnetPair{
				Raw:    s,
				Parsed: parseCIDRs([]string{s}),
			}
			filterCIDRs = append(filterCIDRs, pair)
		case INPUT_TYPE_IP_RANGE:
			cidrs, err := rangeToCidrs(s)
			if err == nil {
				pair := SubnetPair{
					Raw:    s,
					Parsed: parseCIDRs(cidrs),
				}
				filterCIDRs = append(filterCIDRs, pair)
			}
		default:
			return ErrInvalidInput
		}
		return nil
	}

	var err error

	for _, expr := range f.Expression {
		if expr == "-" && isStdin {
			err = ProcessStringsFromStdin(source_op)
			if err != nil {
				return err
			}
		} else {
			err = ProcessStringsFromFile(expr, source_op)
			if err != nil {
				return err
			}
		}
	}

	for _, arg := range args {
		if arg == "-" && isStdin {
			err = ProcessStringsFromStdin(filter_op)
			if err != nil {
				return err
			}
		} else {
			err = ProcessStringsFromFile(arg, filter_op)
			if err != nil {
				return err
			}
		}
	}

	matches := findOverlapping(sourceCIDRs, filterCIDRs, sourceIPs, filterIPs)
	for _, v := range matches {
		fmt.Println(v)
	}

	return nil
}

func rangeToCidrs(r string) ([]string, error) {
	if strings.ContainsRune(r, ':') {
		cidrs, err := CIDRsFromIP6RangeStrRaw(r)
		if err != nil {
			return nil, err
		}
		return cidrs, nil
	} else {
		cidrs, err := CIDRsFromIPRangeStrRaw(r)
		if err != nil {
			return nil, err
		}
		return cidrs, nil
	}
}

type SubnetPair struct {
	Raw    string
	Parsed []net.IPNet
}

func findOverlapping(sourceCIDRs, filterCIDRs []SubnetPair, sourceIPs, filterIPs []net.IP) []string {
	var matches []string
	for _, sourceCIDR := range sourceCIDRs {
		foundMatch := false
		for _, v := range sourceCIDR.Parsed {
			for _, filterCIDR := range filterCIDRs {
				for _, fv := range filterCIDR.Parsed {
					if isCIDROverlapping(&v, &fv) {
						if !foundMatch {
							matches = append(matches, sourceCIDR.Raw)
							foundMatch = true
						}
						break
					}
				}
			}
		}
	}

	for _, sourceIP := range sourceIPs {
		foundMatch := false
		for _, filterCIDR := range filterCIDRs {
			for _, fv := range filterCIDR.Parsed {
				if isIPRangeOverlapping(&sourceIP, &fv) {
					if !foundMatch {
						matches = append(matches, sourceIP.String())
						foundMatch = true
					}
					break
				}
			}
		}

		for _, filterIP := range filterIPs {
			if sourceIP.String() == filterIP.String() {
				matches = append(matches, sourceIP.String())
				break
			}
		}
	}

	return matches
}

func isCIDROverlapping(cidr1, cidr2 *net.IPNet) bool {
	return cidr1.Contains(cidr2.IP) || cidr2.Contains(cidr1.IP)
}

func isIPRangeOverlapping(ip *net.IP, cidr *net.IPNet) bool {
	return cidr.Contains(*ip)
}
