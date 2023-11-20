package lib

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
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

	// Parses a list of CIDRs.
	parseCIDRs := func(cidrs []string) []net.IPNet {
		parsedCIDRs := make([]net.IPNet, 0)
		for _, cidrStr := range cidrs {
			_, ipNet, err := net.ParseCIDR(cidrStr)
			if err != nil {
				if !f.Quiet {
					fmt.Printf("Invalid CIDR: %s\n", cidrStr)
				}
				continue
			}
			parsedCIDRs = append(parsedCIDRs, *ipNet)
		}

		return parsedCIDRs
	}

	// Input parser.
	parseInput := func(rows []string) ([]net.IPNet, []net.IP) {
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
					if !f.Quiet {
						fmt.Printf("Invalid IP range: %s\n", rowStr)
					}
					continue
				}

				if strings.ContainsRune(rowStr, ':') {
					cidrs, err := CIDRsFromIP6RangeStrRaw(rowStr)
					if err == nil {
						parsedCIDRs = append(parsedCIDRs, parseCIDRs(cidrs)...)
						continue
					} else {
						if !f.Quiet {
							fmt.Printf("Invalid IP range %s. Err: %v\n", rowStr, err)
						}
						continue
					}
				} else {
					cidrs, err := CIDRsFromIPRangeStrRaw(rowStr)
					if err == nil {
						parsedCIDRs = append(parsedCIDRs, parseCIDRs(cidrs)...)
						continue
					} else {
						if !f.Quiet {
							fmt.Printf("Invalid IP range %s. Err: %v\n", rowStr, err)
						}
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
	parsedCIDRs := make([]net.IPNet, 0)
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

	// Sort CIDRs by starting IP.
	sortCIDRs(parsedCIDRs)

	// Remove prefixes which are included in another prefix.
	merged := mergeOverlapping(parsedCIDRs)

	// Combine adjacent entries.
	adjacentCombined := combineAdjacent(merged)

	// Print the aggregated CIDRs.
	for _, r := range adjacentCombined {
		fmt.Println(r.String())
	}

	return nil
}

// The adjacency condition is that the prefixes have the same mask length,
// and the second prefix is exactly one larger than the first prefix.
func areAdjacent(r1, r2 net.IPNet) bool {
	prefix1 := binary.BigEndian.Uint32(r1.IP.To4())
	prefix2 := binary.BigEndian.Uint32(r2.IP.To4())

	mask1, _ := r1.Mask.Size()
	mask2, _ := r2.Mask.Size()

	return mask1 == mask2 && (prefix1%(2<<(32-mask1)) == 0) && (prefix2-prefix1 == (1 << (32 - mask1)))
}

func combineAdjacentCIDRs(r1, r2 net.IPNet) net.IPNet {
	mask1, _ := r1.Mask.Size()

	commonPrefixLen := mask1 - 1
	commonPrefix := r1.IP.Mask(r1.Mask)

	return net.IPNet{IP: commonPrefix, Mask: net.CIDRMask(commonPrefixLen, len(commonPrefix)*8)}
}

func combineAdjacent(cidrs []net.IPNet) []net.IPNet {
	res := make([]net.IPNet, 0)

	for i := 0; i < len(cidrs)-1; i++ {

		if areAdjacent(cidrs[i], cidrs[i+1]) {
			res = append(res, combineAdjacentCIDRs(cidrs[i], cidrs[i+1]))
			i++
		} else {
			res = append(res, cidrs[i])
			if i == len(cidrs)-2 {
				res = append(res, cidrs[i+1])
			}
		}
	}

	return res
}

// Helper function to aggregate IP ranges.
func mergeOverlapping(cidrs []net.IPNet) []net.IPNet {
	merged := make([]net.IPNet, 0)

	// Sort CIDRs by starting IP.
	for _, r := range cidrs {
		if len(merged) == 0 {
			merged = append(merged, r)
			continue
		}

		last := len(merged) - 1
		prev := merged[last]

		if canAggregate(prev, r) {
			// Merge overlapping CIDRs.
			merged[last] = merge(prev, r)
		} else {
			merged = append(merged, r)
		}
	}

	return merged
}

// Helper function to sort IP ranges by starting IP.
func sortCIDRs(ipRanges []net.IPNet) {
	sort.SliceStable(ipRanges, func(i, j int) bool {
		return bytes.Compare(ipRanges[i].IP, ipRanges[j].IP) < 0
	})
}

// Helper function to check if two CIDRs can be aggregated.
func canAggregate(r1, r2 net.IPNet) bool {
	return r1.Contains(r2.IP) || r2.Contains(r1.IP)
}

// Helper function to aggregate two CIDRs.
func merge(r1, r2 net.IPNet) net.IPNet {
	mask1, _ := r1.Mask.Size()
	mask2, _ := r2.Mask.Size()

	ipLen := net.IPv6len * 8
	if r1.IP.To4() != nil {
		ipLen = net.IPv4len * 8
	}

	// Find the common prefix length
	commonPrefixLen := mask1
	if mask2 < commonPrefixLen {
		commonPrefixLen = mask2
	}

	commonPrefix := r1.IP.Mask(net.CIDRMask(commonPrefixLen, ipLen))

	return net.IPNet{IP: commonPrefix, Mask: net.CIDRMask(commonPrefixLen, ipLen)}
}
