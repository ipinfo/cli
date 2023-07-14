package lib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/big"
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
				firstIP := ipRange[0]
				lastIP := ipRange[1]

				if strings.ContainsRune(rowStr, ':') {
					cidrs, err := getCIDRFromIP6Range(firstIP, lastIP)
					if err == nil {
						parsedCIDRs = append(parsedCIDRs, parseCIDRs([]string{cidrs})...)
						continue
					} else {
						if !f.Quiet {
							fmt.Printf("Invalid IP range %s. Err: %v\n", rowStr, err)
						}
						continue
					}
				} else {
					cidrs, err := getCIDRFromIP4Range(firstIP, lastIP)
					if err == nil {
						parsedCIDRs = append(parsedCIDRs, parseCIDRs([]string{cidrs})...)
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

	// Sort and merge collected CIDRs and IPs.
	aggregatedCIDRs := aggregateCIDRs(parsedCIDRs)
	outlierIPs := make([]net.IP, 0)
	length := len(aggregatedCIDRs)
	for _, ip := range parsedIPs {
		for i, cidr := range aggregatedCIDRs {
			if cidr.Contains(ip) {
				break
			} else if i == length-1 {
				outlierIPs = append(outlierIPs, ip)
			}
		}
	}

	// Print the aggregated CIDRs.
	for _, r := range aggregatedCIDRs {
		fmt.Println(r.String())
	}

	// Print outliers.
	for _, r := range outlierIPs {
		fmt.Println(r.String())
	}

	return nil
}

// Helper function to get CIDR from an IPv6 IP range.
func getCIDRFromIP6Range(firstIP string, lastIP string) (string, error) {
	startIP := net.ParseIP(strings.TrimSpace(firstIP))
	endIP := net.ParseIP(strings.TrimSpace(lastIP))

	if startIP == nil || endIP == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	// Convert the IP addresses to big.Int
	start := big.NewInt(0).SetBytes(startIP)
	end := big.NewInt(0).SetBytes(endIP)

	if start.Cmp(end) > 0 {
		return "", fmt.Errorf("start IP must be less than or equal to end IP")
	}

	// Calculate the number of host bits required
	numHosts := big.NewInt(0).Sub(end, start)
	numHosts.Add(numHosts, big.NewInt(1))
	numBits := 128 - numHosts.BitLen()

	// Calculate the CIDR notation
	ip := startIP.String()
	cidr := fmt.Sprintf("%s/%d", ip, numBits)

	return cidr, nil
}

// Helper function to get CIDR from an IPv4 IP range.
func getCIDRFromIP4Range(firstIP string, lastIP string) (string, error) {
	startIP := net.ParseIP(strings.TrimSpace(firstIP))
	endIP := net.ParseIP(strings.TrimSpace(lastIP))

	if startIP == nil || endIP == nil {
		return "", fmt.Errorf("invalid IP address")
	}

	// Convert IP addresses to big integers
	start := big.NewInt(0)
	start.SetBytes(startIP.To4())
	end := big.NewInt(0)
	end.SetBytes(endIP.To4())

	// Calculate the subnet mask length
	mask := big.NewInt(0)
	mask.Sub(end, start)
	mask.Add(mask, big.NewInt(1))
	maskLen := mask.BitLen()

	// Format CIDR notation
	cidr := fmt.Sprintf("%s/%d", startIP.String(), 32-maskLen)

	return cidr, nil
}

// Helper function to aggregate IP ranges.
func aggregateCIDRs(cidrs []net.IPNet) []net.IPNet {
	aggregatedCIDRs := make([]net.IPNet, 0)

	// Sort CIDRs by starting IP.
	sortCIDRs(cidrs)

	for _, r := range cidrs {
		if len(aggregatedCIDRs) == 0 {
			aggregatedCIDRs = append(aggregatedCIDRs, r)
			continue
		}

		last := len(aggregatedCIDRs) - 1
		prev := aggregatedCIDRs[last]

		if canAggregate(prev, r) {
			// Merge overlapping CIDRs.
			aggregatedCIDRs[last] = aggregateCIDR(prev, r)
		} else {
			aggregatedCIDRs = append(aggregatedCIDRs, r)
		}
	}

	return aggregatedCIDRs
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
func aggregateCIDR(r1, r2 net.IPNet) net.IPNet {
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
