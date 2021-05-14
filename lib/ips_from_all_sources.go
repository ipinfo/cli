package lib

import (
	"net"
)

// IPsFromAllSources is the same as IPsFrom with all flags turned on.
func IPsFromAllSources(inputs []string) ([]net.IP, error) {
	return IPsFrom(inputs, true, true, true, true, true)
}
