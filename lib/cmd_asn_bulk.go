package lib

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

// CmdASNBulkFlags are flags expected by CmdASNBulk
type CmdASNBulkFlags struct {
	Token   string
	nocache bool
	help    bool
	Field   []string
	json    bool
	Yaml    bool
}

// Init initializes the common flags available to CmdASNBulk with sensible
// defaults.
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdASNBulkFlags) Init() {
	_h := "see description in --help"
	pflag.StringVarP(
		&f.Token,
		"token", "t", "",
		_h,
	)
	pflag.BoolVarP(
		&f.nocache,
		"nocache", "", false,
		_h,
	)
	pflag.BoolVarP(
		&f.help,
		"help", "h", false,
		_h,
	)
	pflag.StringSliceVarP(
		&f.Field,
		"field", "f", []string{},
		_h,
	)
	pflag.BoolVarP(
		&f.json,
		"json", "j", false,
		_h,
	)
	pflag.BoolVarP(
		&f.Yaml,
		"yaml", "y", false,
		_h,
	)
}

// CmdASNBulk is the entrypoint for the `ipinfo asn-bulk` command.
func CmdASNBulk(f CmdASNBulkFlags, ii *ipinfo.Client, args []string, printHelp func()) (ipinfo.BatchASNDetails, error) {
	if f.help {
		printHelp()
		return nil, nil
	}

	var asns []string

	if len(args) == 0 {
		fmt.Println("** manual input mode **\nEnter all ASNs, one per line:")
		args = ReadStringsFromStdin()
		if len(args) == 0 {
			return nil, errors.New("no input ASNs")
		}
	}

	for i := 0; i < len(args); i++ {
		if strings.HasSuffix(args[i], ".txt") {
			lines, err := readStringsFromFile(args[i])
			if err != nil {
				return nil, err
			}
			// Remove arg from args
			args = append(args[:i], args[i+1:]...)
			i-- // Adjust the index as the slice length changes
			args = append(args, lines...)
		}
	}

	if ii.Token == "" {
		return nil, errors.New("bulk lookups require a token; login via `ipinfo init`.")
	}

	for _, arg := range args {
		// Convert to uppercase
		asnUpperCase := strings.ToUpper(arg)
		// Validate ASN
		if !StrIsASNStr(asnUpperCase) {
			return nil, ErrInvalidInput
		}
		asns = append(asns, asnUpperCase)
	}

	return ii.GetASNDetailsBatch(asns, ipinfo.BatchReqOpts{
		TimeoutPerBatch:              60 * 30, // 30min
		ConcurrentBatchRequestsLimit: 20,
	})
}

// ReadStringsFromStdin reads strings from stdin until an empty line is entered.
func ReadStringsFromStdin() []string {
	var inputLines []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		inputLines = append(inputLines, line)
	}
	return inputLines
}

// readStringsFromFile reads strings from a file, one per line.
func readStringsFromFile(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return // ignore errors on close
		}
	}(file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords) // Set the scanner to split on spaces and newlines

	for scanner.Scan() {
		word := scanner.Text()
		lines = append(lines, word)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
