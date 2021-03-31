package lib

import (
	"net"
	"os"
)

// IPsFromStdin returns a list of IPs from a stdin; the IPs should be 1 per
// line.
func IPsFromStdin() []net.IP {
	return IPsFromReader(os.Stdin)
}
