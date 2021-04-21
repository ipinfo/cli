package main

import (
	"fmt"
	"net"
	"os"

	"github.com/ipinfo/cli/lib"
)

func getInputIPs(args []string) ([]net.IP, error) {
	// check for stdin, implied or explicit.
	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Println("** manual input mode **")
			fmt.Println("Enter all IPs, one per line:")
		}
		return lib.IPsFromStdin(), nil
	}

	// check for IP range.
	if lib.IsIP(args[0]) {
		if len(args) != 2 {
			return nil, lib.ErrIPRangeRequiresTwoIPs
		}
		if !lib.IsIP(args[1]) {
			return nil, lib.ErrNotIP
		}

		return lib.IPsFromRange(args[0], args[1])
	}

	// check for all CIDRs.
	if lib.IsCIDR(args[0]) {
		for _, arg := range args[1:] {
			if !lib.IsCIDR(arg) {
				return nil, lib.ErrNotCIDR
			}
		}

		return lib.IPsFromCIDRs(args)
	}

	// check for all filepaths.
	if fileExists(args[0]) {
		for _, arg := range args[1:] {
			if !fileExists(arg) {
				return nil, lib.ErrNotFile
			}
		}

		return lib.IPsFromFiles(args)
	}

	return nil, nil
}
