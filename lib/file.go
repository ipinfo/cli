package lib

import (
	"os"
)

// FileExists checks if a file exists.
func FileExists(pathToFile string) bool {
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		return false
	}
	return true
}
