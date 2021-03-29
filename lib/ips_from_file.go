package lib

import (
	"os"
	"net"
)

func IPsFromFile(pathToFile string) ([]net.IP, error) {
	f, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}

	return IPsFromReader(f), nil
}
