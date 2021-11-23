package main

import (
	"path/filepath"
	"os"
)

// gets the global config directory, creating it if necessary.
func getConfigDir() (string, error) {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	confDir := filepath.Join(cdir, "ipinfo")
	if err := os.MkdirAll(confDir, 0700); err != nil {
		return "", err
	}

	return confDir, nil
}
