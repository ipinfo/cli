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

// IPListFrom returns a list of IPs from stdin and a list of inputs which is
// interpreted to contain IPs, IP ranges, IP CIDRs and files with IPs in them,
// all depending upon which flags are set.
func IPListFrom(
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
			ips = append(ips, IPListFromStdin()...)
		}
	}

	// parse `inputs`.
	for _, input := range inputs {
		if iprange {
			_ips, err := IPListFromRangeStr(input)
			if err == nil {
				ips = append(ips, _ips...)
				continue
			}
		}

		if ip && StrIsIPStr(input) {
			ips = append(ips, net.ParseIP(input))
			continue
		}

		if cidr && StrIsCIDRStr(input) {
			_ips, _ := IPListFromCIDR(input)
			ips = append(ips, _ips...)
			continue
		}

		if file && FileExists(input) {
			_ips, err := IPListFromFile(input)
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

// IPListFromAllSrcs is the same as IPListFrom with all flags turned on.
func IPListFromAllSrcs(inputs []string) ([]net.IP, error) {
	return IPListFrom(inputs, true, true, true, true, true)
}

// IPListFromCIDR returns a list of IPs from a CIDR string.
func IPListFromCIDR(cidrStr string) ([]net.IP, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, err
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)
	end := (start & mask) | (mask ^ 0xffffffff)

	ips := make([]net.IP, 0, end-start+1)
	for i := start; i <= end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip)
	}

	return ips, nil
}

// IPListFromCIDRs returns a list of IPs from a list of CIDRs in string form.
func IPListFromCIDRs(cidrStrs []string) (ips []net.IP, err error) {
	// collect IPs lists together first, then allocate a final list and do
	// a fast transfer.
	ipRanges := make([][]net.IP, len(cidrStrs))
	totalIPs := 0
	for i, cidr := range cidrStrs {
		ipRanges[i], err = IPListFromCIDR(cidr)
		if err != nil {
			return nil, err
		}
		totalIPs += len(ipRanges[i])
	}

	ips = make([]net.IP, 0, totalIPs)
	for _, ipRange := range ipRanges {
		ips = append(ips, ipRange...)
	}

	return ips, nil
}

// IPListFromFile returns a list of IPs found in a file.
func IPListFromFile(pathToFile string) ([]net.IP, error) {
	f, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}

	return IPListFromReader(f), nil
}

// IPListFromFiles returns a list of IPs found in a list of files.
func IPListFromFiles(paths []string) (ips []net.IP, err error) {
	// collect IPs lists together first, then allocate a final list and do
	// a fast transfer.
	ipLists := make([][]net.IP, len(paths))
	totalIPs := 0
	for i, p := range paths {
		ipLists[i], err = IPListFromFile(p)
		if err != nil {
			return nil, err
		}
		totalIPs += len(ipLists[i])
	}

	ips = make([]net.IP, 0, totalIPs)
	for _, ipList := range ipLists {
		ips = append(ips, ipList...)
	}

	return ips, nil
}

// IPListFromRange returns a list of IPs from a start and end IP string.
func IPListFromRange(ipStrStart string, ipStrEnd string) ([]net.IP, error) {
	var ips []net.IP
	var ipStart, ipEnd net.IP

	if ipStart = net.ParseIP(ipStrStart); ipStart == nil {
		return nil, ErrNotIP
	}
	if ipEnd = net.ParseIP(ipStrEnd); ipEnd == nil {
		return nil, ErrNotIP
	}

	start := binary.BigEndian.Uint32(ipStart.To4())
	end := binary.BigEndian.Uint32(ipEnd.To4())
	if start > end {
		ips = make([]net.IP, 0, start-end+1)
		// return decreasing list if range is flipped.
		for i := start; i >= end; i-- {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			ips = append(ips, ip)
		}
	} else {
		ips = make([]net.IP, 0, end-start+1)
		for i := start; i <= end; i++ {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			ips = append(ips, ip)
		}
	}

	return ips, nil
}

// IPListFromRangeStr returns a list of IPs given a range string.
//
// `rStr` must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func IPListFromRangeStr(rStr string) ([]net.IP, error) {
	r, err := IPRangeStrFromStr(rStr)
	if err != nil {
		return nil, err
	}

	return IPListFromRange(r.Start, r.End)
}

// IPListFromReader returns a list of IPs after reading from a reader; the reader
// should have IPs per-line.
func IPListFromReader(r io.Reader) []net.IP {
	ips := make([]net.IP, 0, 10000)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ipStr := strings.TrimSpace(scanner.Text())
		if ipStr == "" {
			break
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			// ignore any non-IP input.
			continue
		}

		ips = append(ips, ip)
	}

	return ips
}

// IPListFromStdin returns a list of IPs from a stdin; the IPs should be 1 per
// line.
func IPListFromStdin() []net.IP {
	return IPListFromReader(os.Stdin)
}

// IPListFromWrite outputs a list of IPs from inputs which are interpreted to
// contain IP ranges and IP CIDRs in them, all depending upon which flags are
// set.
func IPListFromWrite(
	inputs []string,
	iprange bool,
	cidr bool,
) error {
	// prevent edge cases with all flags turned off.
	if !iprange && !cidr {
		return nil
	}

	for _, input := range inputs {
		if iprange {
			if err := IPListFromIPRangeStrWrite(input); err == nil {
				continue
			}
		}

		if cidr && StrIsCIDRStr(input) {
			if err := IPListFromCIDRWrite(input); err == nil {
				continue
			}
		}

		return ErrInvalidInput
	}

	return nil
}

// IPListFromCIDRWrite is the same as IPListFromCIDR with O(1) memory by discarding
// IPs after printing.
func IPListFromCIDRWrite(cidrStr string) error {
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

// IPListFromIPRangeWrite is the same as IPListFromRange with O(1) memory by
// discarding IPs after printing.
func IPListFromIPRangeWrite(ipStrStart string, ipStrEnd string) error {
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
func IPListFromIPRangeStrWrite(rStr string) error {
	r, err := IPRangeStrFromStr(rStr)
	if err != nil {
		return err
	}

	return IPListFromIPRangeWrite(r.Start, r.End)
}
