package lib

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func HelpDetailed(detailedHelp string) error {

	pagerCmd := os.Getenv("PAGER")

	if pagerCmd == "" {
		// If PAGER is not set, use a default pager (e.g., less)
		pagerCmd = "less"
	}

	cmd := exec.Command(pagerCmd)

	// Create an io.Reader from the detailed help string
	reader := io.Reader(strings.NewReader(detailedHelp))
	cmd.Stdin = reader

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// returns an error if the command doesn't execute properly. Otherwise, it returns nil and executes the command
	return cmd.Run()

}
