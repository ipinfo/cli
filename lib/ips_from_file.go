package lib

import (
	"net"
	"os"
)

// IPsFromFile returns a list of IPs found in a file.
func IPsFromFile(pathToFile string) ([]net.IP, error) {
	f, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}

	return IPsFromReader(f), nil
}
