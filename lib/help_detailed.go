package lib

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func HelpDetailed(detailedHelp string, printHelpDefault func()) error {
	pagerCmd := os.Getenv("PAGER")
	if pagerCmd == "" {
		// If PAGER is not set, use a default pager (e.g., less)
		pagerCmd = "less"
	}

	cmd := exec.Command(pagerCmd)
	reader := io.Reader(strings.NewReader(detailedHelp))
	cmd.Stdin = reader
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// If an error occurs running the pager, display the default help
	if err := cmd.Run(); err != nil {
		printHelpDefault()
	}

	return nil
}
