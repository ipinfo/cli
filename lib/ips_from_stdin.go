package lib

import (
	"net"
	"os"
)

func IPsFromStdin() []net.IP {
	return IPsFromReader(os.Stdin)
}
