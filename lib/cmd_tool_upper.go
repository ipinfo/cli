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

type CmdToolUpperFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolUpperFlags) Init() {
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

func CmdToolUpper(
	f CmdToolUpperFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}

	parseCIDRsAndIPs := func(items []string) ([]net.IP, error) {
		parsedIPs := make([]net.IP, 0)
		for _, item := range items {
			if strings.ContainsRune(item, '/') {
				ipRange, err := IPRangeStrFromCIDR(item)
				if err != nil {
					return nil, err
				}
				endIP := net.ParseIP(ipRange.End)
				parsedIPs = append(parsedIPs, endIP)
			} else {
				ip := net.ParseIP(item)
				if ip == nil {
					return nil, fmt.Errorf("invalid input: %q", item)
				}
				parsedIPs = append(parsedIPs, ip)
			}
		}
		return parsedIPs, nil
	}
	parsedIPs := make([]net.IP, 0)

	if isStdin {
		rows := scanrdr(os.Stdin)
		ips, err := parseCIDRsAndIPs(rows)
		if err != nil {
			if !f.Quiet {
				fmt.Println(err)
			}
			return nil
		}
		parsedIPs = append(parsedIPs, ips...)
	}

	for _, arg := range args {
		file, err := os.Open(arg)
		if err != nil {
			ips, err := parseCIDRsAndIPs([]string{arg})
			if err != nil {
				if !f.Quiet {
					fmt.Println(err)
				}
			}
			parsedIPs = append(parsedIPs, ips...)
			continue
		}

		rows := scanrdr(file)
		file.Close()
		ips, err := parseCIDRsAndIPs(rows)
		if err != nil {
			if !f.Quiet {
				fmt.Println(err)
			}
		}
		parsedIPs = append(parsedIPs, ips...)
	}

	for _, ip := range parsedIPs {
		fmt.Println(ip.String())
	}

	return nil
}

func scanrdr(r io.Reader) []string {
	rows := make([]string, 0)

	buf := bufio.NewReader(r)
	for {
		d, err := buf.ReadString('\n')
		if err == io.EOF {
			if len(d) == 0 {
				break
			}
		} else if err != nil {
			return rows
		}

		sepIdx := strings.IndexAny(d, "\n")
		if sepIdx == -1 {
			sepIdx = len(d)
		}

		rowStr := d[:sepIdx]
		rows = append(rows, rowStr)
	}

	return rows
}
