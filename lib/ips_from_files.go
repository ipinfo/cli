package lib

import (
	"net"
)

// IPsFromFiles returns a list of IPs found in a list of files.
func IPsFromFiles(paths []string) (ips []net.IP, err error) {
	// collect IPs lists together first, then allocate a final list and do
	// a fast transfer.
	ipLists := make([][]net.IP, len(paths))
	totalIPs := 0
	for i, p := range paths {
		ipLists[i], err = IPsFromFile(p)
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
