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

// CmdToolAggregateFlags are flags expected by CmdToolAggregate.
type CmdToolAggregateFlags struct {
	Help  bool
	Quiet bool
}

// Init initializes the common flags available to CmdToolAggregate with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdToolAggregateFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.Quiet,
		"quiet", "q", false,
		"quiet mode; suppress additional output.",
	)
}

// CmdToolAggregate is the common core logic for aggregating IPs, IP ranges and CIDRs.
func CmdToolAggregate(
	f CmdToolAggregateFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	// require args.
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}

	// Input parser.
	parseInput := func(rows []string) ([]string, []net.IP) {
		parsedCIDRs := make([]string, 0)
		parsedIPs := make([]net.IP, 0)
		for _, rowStr := range rows {
			if strings.ContainsAny(rowStr, ",-") {
				continue
			} else if strings.ContainsRune(rowStr, '/') {
				_, ipnet, err := net.ParseCIDR(rowStr)
				if err == nil && IsCIDRIPv4(ipnet) {
					parsedCIDRs = append(parsedCIDRs, []string{rowStr}...)
				}
				continue
			} else {
				if ip := net.ParseIP(rowStr); IsIPv4(ip) {
					parsedIPs = append(parsedIPs, ip)
				} else {
					if !f.Quiet {
						fmt.Printf("Invalid input: %s\n", rowStr)
					}
				}
			}
		}

		return parsedCIDRs, parsedIPs
	}

	// Input scanner.
	scanrdr := func(r io.Reader) []string {
		rows := make([]string, 0)

		buf := bufio.NewReader(r)
		for {
			d, err := buf.ReadString('\n')
			if err == io.EOF {
				if len(d) == 0 {
					break
				}
			} else if err != nil {
				if !f.Quiet {
					fmt.Printf("Scan error: %v\n", err)
				}
				return rows
			}

			sepIdx := strings.IndexAny(d, "\n")
			if sepIdx == -1 {
				// only possible if EOF & input doesn't end with newline.
				sepIdx = len(d)
			}

			rowStr := d[:sepIdx]
			rows = append(rows, rowStr)
		}

		return rows
	}

	// Vars to contain CIDRs/IPs from all input sources.
	parsedCIDRs := make([]string, 0)
	parsedIPs := make([]net.IP, 0)

	// Collect CIDRs/IPs from stdin.
	if isStdin {
		rows := scanrdr(os.Stdin)
		parsedCIDRs, parsedIPs = parseInput(rows)
	}

	// Collect CIDRs/IPs from all args.
	for _, arg := range args {
		file, err := os.Open(arg)
		if err != nil {
			cidrs, ips := parseInput([]string{arg})
			parsedCIDRs = append(parsedCIDRs, cidrs...)
			parsedIPs = append(parsedIPs, ips...)
			continue
		}

		rows := scanrdr(file)
		file.Close()
		cidrs, ips := parseInput(rows)

		parsedCIDRs = append(parsedCIDRs, cidrs...)
		parsedIPs = append(parsedIPs, ips...)
	}

	adjacentCombined := combineAdjacent(stripOverlapping(list(parsedCIDRs)))

	outlierIPs := make([]net.IP, 0)
	length := len(adjacentCombined)
	if length != 0 {
		for _, ip := range parsedIPs {
			for i, cidr := range adjacentCombined {
				if cidr.Network.Contains(ip) {
					break
				} else if i == length-1 {
					outlierIPs = append(outlierIPs, ip)
				}
			}
		}
	} else {
		outlierIPs = append(outlierIPs, parsedIPs...)
	}

	// Print the aggregated CIDRs.
	for _, r := range adjacentCombined {
		fmt.Println(r.String())
	}

	// Print the outlierIPs.
	for _, r := range outlierIPs {
		fmt.Println(r.String())
	}

	return nil
}

// stripOverlapping returns a slice of CIDR structures with overlapping ranges
// stripped.
func stripOverlapping(s []*CIDR) []*CIDR {
	l := len(s)
	for i := 0; i < l-1; i++ {
		if s[i] == nil {
			continue
		}
		for j := i + 1; j < l; j++ {
			if overlaps(s[j], s[i]) {
				s[j] = nil
			}
		}
	}
	return filter(s)
}

func overlaps(a, b *CIDR) bool {
	return (a.PrefixUint32() / (1 << (32 - b.MaskLen()))) ==
		(b.PrefixUint32() / (1 << (32 - b.MaskLen())))
}

// combineAdjacent returns a slice of CIDR structures with adjacent ranges
// combined.
func combineAdjacent(s []*CIDR) []*CIDR {
	for {
		found := false
		l := len(s)
		for i := 0; i < l-1; i++ {
			if s[i] == nil {
				continue
			}
			for j := i + 1; j < l; j++ {
				if s[j] == nil {
					continue
				}
				if adjacent(s[i], s[j]) {
					c := fmt.Sprintf("%s/%d", s[i].IP.String(), s[i].MaskLen()-1)
					s[i] = newCidr(c)
					s[j] = nil
					found = true
				}
			}
		}

		if !found {
			break
		}
	}
	return filter(s)
}

func adjacent(a, b *CIDR) bool {
	return (a.MaskLen() == b.MaskLen()) &&
		(a.PrefixUint32()%(2<<(32-b.MaskLen())) == 0) &&
		(b.PrefixUint32()-a.PrefixUint32() == (1 << (32 - a.MaskLen())))
}

func filter(s []*CIDR) []*CIDR {
	out := s[:0]
	for _, x := range s {
		if x != nil {
			out = append(out, x)
		}
	}
	return out
}
