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
	INPUT_TYPE_ASN      INPUT_TYPE = "ASN"
	INPUT_TYPE_UNKNOWN  INPUT_TYPE = "Unknown"
)

func inputHelper(str string, op func(string, INPUT_TYPE) error) error {
	switch {
	case StrIsIPStr(str):
		return op(str, INPUT_TYPE_IP)
	case StrIsIPRangeStr(str):
		return op(str, INPUT_TYPE_IP_RANGE)
	case StrIsCIDRStr(str):
		return op(str, INPUT_TYPE_CIDR)
	case StrIsASNStr(str):
		return op(str, INPUT_TYPE_ASN)
	default:
		return op(str, INPUT_TYPE_UNKNOWN)
	}
}

// GetInputFrom retrieves input data from various sources and processes it using the provided operation.
// The operation is called for each input string with input type.
//
//	Usage:
//	err := GetInputFrom(inputs,
//		true,
//		true,
//		func(input string, inputType INPUT_TYPE) error {
//			switch inputType {
//			case INPUT_TYPE_IP:
//				// Process IP here
//			}
//			return nil
//		},
//	)
func GetInputFrom(
	inputs []string,
	stdin bool,
	file bool,
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
			err := ProcessStringsFromStdin(op)
			if err != nil {
				return err
			}
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
		case file && FileExists(input):
			err = ProcessStringsFromFile(input, op)
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

// ProcessStringsFromFile reads strings from a file and passes it to op, one per line.
func ProcessStringsFromFile(filename string, op func(input string, inputType INPUT_TYPE) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
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
		err = inputHelper(scanner.Text(), op)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// ProcessStringsFromStdin reads strings from stdin until an empty line is entered.
func ProcessStringsFromStdin(op func(input string, inputType INPUT_TYPE) error) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		err := inputHelper(line, op)
		if err != nil {
			return err
		}
	}
	return nil
}
