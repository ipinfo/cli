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
	FilterFile   []string
	CriteriaFile []string
	Help         bool
}

func (f *CmdMatchIPFlags) Init() {
	pflag.StringSliceVarP(
		&f.FilterFile,
		"filter", "f", nil,
		"IPs, subnets to be filtered.",
	)
	pflag.StringSliceVarP(
		&f.CriteriaFile,
		"criteria", "c", nil,
		"subnets to check overlap with.",
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
	if f.Help || f.FilterFile[0] == "" || f.CriteriaFile[0] == "" {
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
	if len(f.FilterFile) == 1 && f.FilterFile[0] == "-" && isStdin {
		filter, err = scanrdr(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		for _, file := range f.FilterFile {
			f, err := os.Open(file)
			if err != nil {
				return err
			}

			res, err := scanrdr(f)
			if err != nil {
				return err
			}
			filter = append(filter, res...)
		}
	}

	var criteria []string
	if len(f.CriteriaFile) == 1 && f.CriteriaFile[0] == "-" && isStdin {
		criteria, err = scanrdr(os.Stdin)
		if err != nil {
			return err
		}
	} else {
		for _, file := range f.CriteriaFile {
			f, err := os.Open(file)
			if err != nil {
				return err
			}

			res, err := scanrdr(f)
			if err != nil {
				return err
			}
			criteria = append(criteria, res...)
		}
	}

	matches := findOverlapping(filter, criteria)
	for _, v := range matches {
		fmt.Println(v)
	}

	return nil
}

type SubnetPair struct {
	Raw    string
	Parsed []net.IPNet
}

func findOverlapping(filter, criteria []string) []string {
	parseInput := func(rows []string) ([]SubnetPair, []net.IP) {
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

		parsedCIDRs := make([]SubnetPair, 0)
		parsedIPs := make([]net.IP, 0)
		var separator string
		for _, rowStr := range rows {
			if strings.ContainsAny(rowStr, ",-") {
				if delim := strings.ContainsRune(rowStr, ','); delim {
					separator = ","
				} else {
					separator = "-"
				}

				ipRange := strings.Split(rowStr, separator)
				if len(ipRange) != 2 {
					continue
				}

				if strings.ContainsRune(rowStr, ':') {
					cidrs, err := CIDRsFromIP6RangeStrRaw(rowStr)
					if err == nil {
						pair := SubnetPair{
							Raw:    rowStr,
							Parsed: parseCIDRs(cidrs),
						}
						parsedCIDRs = append(parsedCIDRs, pair)
					} else {
						continue
					}
				} else {
					cidrs, err := CIDRsFromIPRangeStrRaw(rowStr)
					if err == nil {
						pair := SubnetPair{
							Raw:    rowStr,
							Parsed: parseCIDRs(cidrs),
						}
						parsedCIDRs = append(parsedCIDRs, pair)
					} else {
						continue
					}
				}
			} else if strings.ContainsRune(rowStr, '/') {
				pair := SubnetPair{
					Raw:    rowStr,
					Parsed: parseCIDRs([]string{rowStr}),
				}
				parsedCIDRs = append(parsedCIDRs, pair)
			} else {
				if ip := net.ParseIP(rowStr); ip != nil {
					parsedIPs = append(parsedIPs, ip)
				}
			}
		}

		return parsedCIDRs, parsedIPs
	}

	sourceCIDRs, sourceIPs := parseInput(filter)
	filterCIDRs, filterIPs := parseInput(criteria)

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
