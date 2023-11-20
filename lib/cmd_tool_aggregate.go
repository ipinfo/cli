package lib

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
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

// CIDR represens a Classless Inter-Domain Routing structure.
type CIDR struct {
	IP      net.IP
	Network *net.IPNet
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
						parsedCIDRs = append(parsedCIDRs, cidrs...)
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
						parsedCIDRs = append(parsedCIDRs, cidrs...)
						continue
					} else {
						if !f.Quiet {
							fmt.Printf("Invalid IP range %s. Err: %v\n", rowStr, err)
						}
						continue
					}
				}
			} else if strings.ContainsRune(rowStr, '/') {
				parsedCIDRs = append(parsedCIDRs, []string{rowStr}...)
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

	adjacentCombined := CombineAdjacent(StripOverlapping(List(parsedCIDRs)))

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

// New creates a new CIDR structure.
func New(s string) *CIDR {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return &CIDR{
		IP:      ip,
		Network: ipnet,
	}
}

func (c *CIDR) String() string {
	return c.Network.String()
}

// MaskLen returns a network mask length.
func (c *CIDR) MaskLen() uint32 {
	i, _ := c.Network.Mask.Size()
	return uint32(i)
}

// PrefixUint32 returns a prefix.
func (c *CIDR) PrefixUint32() uint32 {
	return binary.BigEndian.Uint32(c.IP.To4())
}

// Size returns a size of a CIDR range.
func (c *CIDR) Size() int {
	ones, bits := c.Network.Mask.Size()
	return int(math.Pow(2, float64(bits-ones)))
}

// List returns a slice of sorted CIDR structures.
func List(s []string) []*CIDR {
	out := make([]*CIDR, 0)
	for _, c := range s {
		out = append(out, New(c))
	}
	sort.Sort(cidrSort(out))
	return out
}

type cidrSort []*CIDR

func (s cidrSort) Len() int      { return len(s) }
func (s cidrSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s cidrSort) Less(i, j int) bool {
	cmp := bytes.Compare(s[i].IP, s[j].IP)
	return cmp < 0 || (cmp == 0 && s[i].MaskLen() < s[j].MaskLen())
}

// StripOverlapping returns a slice of CIDR structures with overlapping ranges
// stripped.
func StripOverlapping(s []*CIDR) []*CIDR {
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

// CombineAdjacent returns a slice of CIDR structures with adjacent ranges
// combined.
func CombineAdjacent(s []*CIDR) []*CIDR {
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
					s[i] = New(c)
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

// Aggregate returns a slice of CIDR structures with adjacent ranges combined
// and overlapping ranges stripped.
func Aggregate(s []string) []*CIDR {
	return CombineAdjacent(StripOverlapping(List(s)))
}

// Size calculates an overal size of a slice of CIDR structures.
func Size(s []*CIDR) (i int) {
	for _, x := range s {
		i += x.Size()
	}
	return
}
