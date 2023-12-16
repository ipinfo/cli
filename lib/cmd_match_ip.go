package lib

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ipinfo/cli/lib/iputil"
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

type SubnetPair struct {
	Raw    string
	Parsed []net.IPNet
}

func CmdMatchIP(
	f CmdMatchIPFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help || (len(f.Expression) == 0 && len(args) < 2) || len(args) == 0 {
		printHelp()
		return nil
	}

	if len(f.Expression) == 0 {
		f.Expression = []string{args[0]}
		args = args[1:]
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

	processInput := func(container *[]SubnetPair, ips *[]net.IP, s string, inputType iputil.INPUT_TYPE) error {
		switch inputType {
		case iputil.INPUT_TYPE_IP:
			*ips = append(*ips, net.ParseIP(s))
		case iputil.INPUT_TYPE_CIDR:
			pair := SubnetPair{
				Raw:    s,
				Parsed: parseCIDRs([]string{s}),
			}
			*container = append(*container, pair)
		case iputil.INPUT_TYPE_IP_RANGE:
			cidrs, err := rangeToCidrs(s)
			if err == nil {
				pair := SubnetPair{
					Raw:    s,
					Parsed: parseCIDRs(cidrs),
				}
				*container = append(*container, pair)
			}
		default:
			return iputil.ErrInvalidInput
		}
		return nil
	}

	source_op := func(s string, inputType iputil.INPUT_TYPE) error {
		return processInput(&sourceCIDRs, &sourceIPs, s, inputType)
	}

	filter_op := func(s string, inputType iputil.INPUT_TYPE) error {
		return processInput(&filterCIDRs, &filterIPs, s, inputType)
	}

	var err error

	for _, expr := range f.Expression {
		if expr == "-" && isStdin {
			err = iputil.ProcessStringsFromStdin(filter_op)
			if err != nil {
				return err
			}
		} else if iputil.FileExists(expr) {
			err = iputil.ProcessStringsFromFile(expr, filter_op)
			if err != nil {
				return err
			}
		} else {
			err = iputil.InputHelper(expr, filter_op)
			if err != nil {
				return err
			}
		}
	}

	for _, arg := range args {
		if arg == "-" && isStdin {
			err = iputil.ProcessStringsFromStdin(source_op)
			if err != nil {
				return err
			}
		} else if iputil.FileExists(arg) {
			err = iputil.ProcessStringsFromFile(arg, source_op)
			if err != nil {
				return err
			}
		} else {
			err = iputil.InputHelper(arg, source_op)
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
		cidrs, err := iputil.CIDRsFromIP6RangeStrRaw(r)
		if err != nil {
			return nil, err
		}
		return cidrs, nil
	} else {
		cidrs, err := iputil.CIDRsFromIPRangeStrRaw(r)
		if err != nil {
			return nil, err
		}
		return cidrs, nil
	}
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
