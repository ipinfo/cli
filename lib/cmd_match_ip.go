package lib

import (
	"bufio"
	"fmt"
	"io"
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

	scanrdr := func(r io.Reader) ([]string, error) {
		var hitEOF bool
		buf := bufio.NewReader(r)
		var ips []string

		for {
			if hitEOF {
				return ips, nil
			}

			d, err := buf.ReadString('\n')
			d = strings.TrimRight(d, "\n")
			if err == io.EOF {
				if len(d) == 0 {
					return ips, nil
				}

				hitEOF = true
			} else if err != nil {
				return ips, err
			}

			if len(d) == 0 {
				continue
			}

			ips = append(ips, d)
		}
	}

	var filter []string
	var err error

	for _, expr := range f.Expression {
		if expr == "-" && isStdin {
			exprs, err := scanrdr(os.Stdin)
			if err != nil {
				return err
			}
			filter = append(filter, exprs...)
		} else {
			file, err := os.Open(expr)
			if err != nil {
				return err
			}

			exprs, err := scanrdr(file)
			if err != nil {
				return err
			}
			filter = append(filter, exprs...)
		}
	}

	var criteria []string
	for _, arg := range args {
		if arg == "-" && isStdin {
			criteria, err = scanrdr(os.Stdin)
			if err != nil {
				return err
			}
		} else {
			file, err := os.Open(arg)
			if err != nil {
				return err
			}

			res, err := scanrdr(file)
			if err != nil {
				return err
			}
			criteria = append(criteria, res...)
		}
	}

	matches := findOverlapping(sourceCIDRs, filterCIDRs, sourceIPs, filterIPs)
	for _, v := range matches {
		fmt.Println(v)
	}

	return nil
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
