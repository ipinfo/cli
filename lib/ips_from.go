package lib

import (
	"fmt"
	"net"
	"os"
)

// IPsFrom returns a list of IPs from stdin and a list of inputs which is
// interpreted to contain IPs, IP ranges, IP CIDRs and files with IPs in them,
// all depending upon which flags are set.
func IPsFrom(
	inputs []string,
	stdin bool,
	ip bool,
	iprange bool,
	cidr bool,
	file bool,
) ([]net.IP, error) {
	ips := make([]net.IP, 0, len(inputs))

	// prevent edge cases with all flags turned off.
	if !stdin && !ip && !iprange && !cidr && !file {
		return ips, nil
	}

	// start with stdin.
	if stdin {
		stat, _ := os.Stdin.Stat()

		isPiped := (stat.Mode() & os.ModeNamedPipe) != 0
		isTyping := (stat.Mode()&os.ModeCharDevice) != 0 && len(inputs) == 0

		if isTyping {
			fmt.Println("** manual input mode **")
			fmt.Println("Enter all IPs, one per line:")
		}

		if isPiped || isTyping || stat.Size() > 0 {
			ips = append(ips, IPsFromStdin()...)
		}
	}

	// parse `inputs`.
	for _, input := range inputs {
		if iprange {
			_ips, err := IPsFromRangeStr(input)
			if err == nil {
				ips = append(ips, _ips...)
				continue
			}
		}

		if ip && IsIP(input) {
			ips = append(ips, net.ParseIP(input))
			continue
		}

		if cidr && IsCIDR(input) {
			_ips, _ := IPsFromCIDR(input)
			ips = append(ips, _ips...)
			continue
		}

		if file && FileExists(input) {
			_ips, err := IPsFromFile(input)
			if err != nil {
				return nil, err
			}
			ips = append(ips, _ips...)
			continue
		}

		return nil, ErrInvalidInput
	}

	return ips, nil
}
