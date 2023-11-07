package lib

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type CmdMatchIPFlags struct {
	FilterFile       string
	OverlapCheckFile string
	Help             bool
}

func (f *CmdMatchIPFlags) Init() {
	pflag.StringVarP(
		&f.FilterFile,
		"filter", "f", "",
		"file containing a list of IPs, CIDRs, and/or Ranges for filtering.",
	)
	pflag.StringVarP(
		&f.OverlapCheckFile,
		"overlap-check", "o", "",
		"file containing a list of CIDRs and/or Ranges to check for overlap.",
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
	if f.Help || f.FilterFile == "" || f.OverlapCheckFile == "" {
		printHelp()
		return nil
	}

	fmt.Println(f.FilterFile, f.OverlapCheckFile)

	source, err := scanrdr(f.FilterFile)
	if err != nil {
		fmt.Printf("Error reading filter file: %v\n", err)
		return err
	}

	filter, err := scanrdr(f.OverlapCheckFile)
	if err != nil {
		fmt.Printf("Error reading overlap check file: %v\n", err)
		return err
	}

	matches := findOverlappingIPs(source, filter)

	fmt.Println("matches:", matches)

	return nil
}

func scanrdr(filename string) ([]string, error) {
	var ips []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ips = append(ips, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ips, nil
}

func parseInput(rows []string) ([]net.IPNet, []net.IP) {
	parseCIDRs := func(cidrs []string) []net.IPNet {
		parsedCIDRs := make([]net.IPNet, 0)
		for _, cidrStr := range cidrs {
			_, ipNet, err := net.ParseCIDR(cidrStr)
			if err != nil {
				fmt.Printf("Invalid CIDR: %s\n", cidrStr)
				continue
			}
			parsedCIDRs = append(parsedCIDRs, *ipNet)
		}

		return parsedCIDRs
	}

	parsedCIDRs := make([]net.IPNet, 0)
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
				fmt.Printf("Invalid IP range: %s\n", rowStr)
				continue
			}

			if strings.ContainsRune(rowStr, ':') {
				cidrs, err := CIDRsFromIP6RangeStrRaw(rowStr)
				if err == nil {
					parsedCIDRs = append(parsedCIDRs, parseCIDRs(cidrs)...)
					continue
				} else {
					fmt.Printf("Invalid IP range %s. Err: %v\n", rowStr, err)
					continue
				}
			} else {
				cidrs, err := CIDRsFromIPRangeStrRaw(rowStr)
				if err == nil {
					parsedCIDRs = append(parsedCIDRs, parseCIDRs(cidrs)...)
					continue
				} else {
					fmt.Printf("Invalid IP range %s. Err: %v\n", rowStr, err)
					continue
				}
			}
		} else if strings.ContainsRune(rowStr, '/') {
			parsedCIDRs = append(parsedCIDRs, parseCIDRs([]string{rowStr})...)
			continue
		} else {
			if ip := net.ParseIP(rowStr); ip != nil {
				parsedIPs = append(parsedIPs, ip)
			} else {
				fmt.Printf("Invalid input: %s\n", rowStr)
			}
		}
	}

	return parsedCIDRs, parsedIPs
}

func isCIDROverlapping(cidr1, cidr2 *net.IPNet) bool {
	return cidr1.Contains(cidr2.IP) || cidr2.Contains(cidr1.IP)
}

func isIPRangeOverlapping(ip *net.IP, cidr *net.IPNet) bool {
	return cidr.Contains(*ip)
}

func findOverlappingIPs(source, filter []string) []string {
	var matches []string

	sourceCIDRs, sourceIPs := parseInput(source)
	filterCIDRs, filterIPs := parseInput(filter)

	for _, sourceCIDR := range sourceCIDRs {
		for _, filterCIDR := range filterCIDRs {
			if isCIDROverlapping(&sourceCIDR, &filterCIDR) {
				matches = append(matches, sourceCIDR.String())
				break
			}
		}
	}

	for _, sourceIP := range sourceIPs {
		for _, filterCIDR := range filterCIDRs {
			if isIPRangeOverlapping(&sourceIP, &filterCIDR) {
				matches = append(matches, sourceIP.String())
				break
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
