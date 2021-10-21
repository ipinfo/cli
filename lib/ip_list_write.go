package lib

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// IPListFromWrite outputs a list of IPs from stdin and a list of inputs which
// are interpreted to contain IPs, IP ranges, IP CIDRs and files with IPs in
// them, all depending upon which flags are set.
func IPListWriteFrom(
	inputs []string,
	stdin bool,
	ip bool,
	iprange bool,
	cidr bool,
	file bool,
) error {
	// prevent edge cases with all flags turned off.
	if !stdin && !ip && !iprange && !cidr && !file {
		return nil
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
			IPListWriteFromStdin()
		}
	}

	// parse `inputs`.
	for _, input := range inputs {
		if iprange {
			if err := IPListWriteFromIPRangeStr(input); err == nil {
				continue
			}
		}

		if ip && StrIsIPStr(input) {
			fmt.Println(input)
			continue
		}

		if cidr && StrIsCIDRStr(input) {
			if err := IPListWriteFromCIDR(input); err == nil {
				continue
			}
		}

		if file && FileExists(input) {
			if err := IPListWriteFromFile(input); err == nil {
				continue
			}
		}

		return ErrInvalidInput
	}

	return nil
}

// IPListFromAllSrcsWrite is the same as IPListFromWrite with all flags turned
// on.
func IPListWriteFromAllSrcs(inputs []string) error {
	return IPListWriteFrom(inputs, true, true, true, true, true)
}

// IPListFromCIDRWrite is the same as IPListFromCIDR with O(1) memory by discarding
// IPs after printing.
func IPListWriteFromCIDR(cidrStr string) error {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return err
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)
	end := (start & mask) | (mask ^ 0xffffffff)

	for i := start; i <= end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		fmt.Println(ip)
	}

	return nil
}

// IPListFromCIDRsWrite outputs a list of IPs from a list of CIDRs in string
// form.
func IPListWriteFromCIDRs(cidrStrs []string) error {
	for _, cidr := range cidrStrs {
		if err := IPListWriteFromCIDR(cidr); err != nil {
			return err
		}
	}
	return nil
}

// IPListFromIPRangeWrite is the same as IPListFromRange with O(1) memory by
// discarding IPs after printing.
func IPListWriteFromIPRange(ipStrStart string, ipStrEnd string) error {
	var ipStart, ipEnd net.IP

	if ipStart = net.ParseIP(ipStrStart); ipStart == nil {
		return ErrNotIP
	}
	if ipEnd = net.ParseIP(ipStrEnd); ipEnd == nil {
		return ErrNotIP
	}

	start := binary.BigEndian.Uint32(ipStart.To4())
	end := binary.BigEndian.Uint32(ipEnd.To4())

	if start > end {
		// return decreasing list if range is flipped.
		for i := start; i >= end; i-- {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			fmt.Println(ip)
		}
	} else {
		for i := start; i <= end; i++ {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			fmt.Println(ip)
		}
	}

	return nil
}

// IPListFromIPRangeStrWrite outputs all IPs in an IP range.
//
// `rStr` must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func IPListWriteFromIPRangeStr(rStr string) error {
	r, err := IPRangeStrFromStr(rStr)
	if err != nil {
		return err
	}

	return IPListWriteFromIPRange(r.Start, r.End)
}

// IPListWriteFromReader returns a list of IPs after reading from a reader; the
// reader should have IPs per-line.
func IPListWriteFromReader(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ipStr := strings.TrimSpace(scanner.Text())
		if ipStr == "" {
			break
		}

		if err := IPListWriteFromIPRangeStr(ipStr); err == nil {
			continue
		}

		if StrIsIPStr(ipStr) {
			fmt.Println(ipStr)
			continue
		}

		if StrIsCIDRStr(ipStr) {
			if err := IPListWriteFromCIDR(ipStr); err == nil {
				continue
			}
		}

		// simply ignore anything else.
	}
}

// IPListWriteFromStdin returns a list of IPs from a stdin; the IPs should be 1
// per line.
func IPListWriteFromStdin() {
	IPListWriteFromReader(os.Stdin)
}

// IPListWriteFromFile returns a list of IPs found in a file.
func IPListWriteFromFile(pathToFile string) error {
	f, err := os.Open(pathToFile)
	if err != nil {
		return err
	}

	IPListWriteFromReader(f)
	return nil
}

// IPListWriteFromFiles returns a list of IPs found in a list of files.
func IPListWriteFromFiles(paths []string) error {
	for _, p := range paths {
		if err := IPListWriteFromFile(p); err != nil {
			return err
		}
	}
	return nil
}
