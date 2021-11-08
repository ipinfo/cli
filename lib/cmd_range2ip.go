package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// CmdRange2IPFlags are flags expected by CmdRange2IP.
type CmdRange2IPFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdRange2IP with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdRange2IPFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}
func CmdRange2IP(f CmdRange2IPFlags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	// require args and/or stdin.
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}
	// if gets piped input
	isPiped := (stat.Mode() & os.ModeNamedPipe) != 0
	if isPiped || stat.Size() > 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ipStr := strings.TrimSpace(scanner.Text())
			if ipStr == "" {
				break
			}

			if err := IPListWriteFromIPRangeStr(ipStr); err == nil {
				continue
			}

			if StrIsIPStr(ipStr) {
				fmt.Println(ipStr)
				continue
			}

		}
	}
	// reading input
	for _, input := range args {
		f, err := os.Open(input)
		if err != nil {
			// if input is ip
			if StrIsIPStr(input) {
				fmt.Println(input)
				continue
			}
			// if input is ip range
			if err := IPListWriteFromIPRangeStr(input); err == nil {
				continue
			}
			return err
		}
		// if input is file
		if FileExists(input) {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				ipStr := strings.TrimSpace(scanner.Text())
				if ipStr == "" {
					break
				}
				err := IPListWriteFromIPRangeStr(ipStr)
				if err != nil {
					return ErrInvalidInput
				}
				if StrIsIPStr(ipStr) {
					fmt.Println(ipStr)
					continue
				}
			}
		}

	}
	return nil
}
