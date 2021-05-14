package lib

import (
	"bufio"
	"io"
	"net"
	"strings"
)

// IPsFromReader returns a list of IPs after reading from a reader; the reader
// should have IPs per-line.
func IPsFromReader(r io.Reader) []net.IP {
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
