package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type INPUT_TYPE string

const (
	INPUT_TYPE_IP       INPUT_TYPE = "IP"
	INPUT_TYPE_IP_RANGE INPUT_TYPE = "IPRange"
	INPUT_TYPE_CIDR     INPUT_TYPE = "CIDR"
	INPUT_TYPE_FILE     INPUT_TYPE = "File"
	INPUT_TYPE_ASN      INPUT_TYPE = "ASN"
	INPUT_TYPE_UNKNOWN  INPUT_TYPE = "Unknown"
)

// getInputFrom retrieves input data from various sources and processes it using the provided operation.
// The operation is called for each input string with input type.
//
//	Usage:
//	err := getInputFrom(inputs,
//		true,
//		func(input string, inputType INPUT_TYPE) error {
//			switch inputType {
//			case INPUT_TYPE_IP:
//				// Process IP here
//			default:
//				return ErrNotIP
//			}
//			return nil
//		})
func getInputFrom(
	inputs []string,
	stdin bool,
	op func(input string, inputType INPUT_TYPE) error,
) error {
	if !stdin && len(inputs) == 0 {
		return nil
	}

	// start with stdin.
	if stdin {
		stat, _ := os.Stdin.Stat()

		isPiped := (stat.Mode() & os.ModeNamedPipe) != 0
		isTyping := (stat.Mode()&os.ModeCharDevice) != 0 && len(inputs) == 0

		if isTyping {
			fmt.Println("** manual input mode **")
			fmt.Println("one input per line:")
		}

		if isPiped || isTyping || stat.Size() > 0 {
			inputs = append(inputs, ReadStringsFromStdin()...)
		}
	}

	// parse `inputs`.
	for _, input := range inputs {
		var err error
		switch {
		case StrIsIPStr(input):
			err = op(input, INPUT_TYPE_IP)
		case StrIsIPRangeStr(input):
			err = op(input, INPUT_TYPE_IP_RANGE)
		case StrIsCIDRStr(input):
			err = op(input, INPUT_TYPE_CIDR)
		case FileExists(input):
			err = op(input, INPUT_TYPE_FILE)
		case StrIsASNStr(input):
			err = op(input, INPUT_TYPE_ASN)
		default:
			err = op(input, INPUT_TYPE_UNKNOWN)
		}
		if err != nil {
			return err
		}
	}
	return nil
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
