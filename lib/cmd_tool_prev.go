package lib

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type CmdToolPrevFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolPrevFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.Quiet,
		"quiet", "q", false,
		"quiet mode; suppress additional output.",
	)
}

func CmdToolPrev(
	f CmdToolPrevFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	actionStdin := func(input string) {
		ActionForStdinPrev(input)
	}
	actionFile := func(input string) {
		ActionForFilePrev(input)
	}

	// Process inputs using the IPInputAction function.
	err := IPInputAction(args, true, true, false, false, true,
		actionStdin, nil, nil, actionFile)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionForStdinPrev(input string) {
	ip := net.ParseIP(input)
	if ip != nil {
		prevIP := PrevIP(ip)
		fmt.Println(prevIP)
	}
}

func ActionForFilePrev(pathToFile string) {
	f, err := os.Open(pathToFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		ActionForStdinPrev(input)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func PrevIP(ip net.IP) net.IP {
	prevIP := make(net.IP, len(ip))
	copy(prevIP, ip)

	for i := len(prevIP) - 1; i >= 0; i-- {
		if prevIP[i] > 0 {
			prevIP[i]--
			break
		} else {
			prevIP[i] = 255
		}
	}

	return prevIP
}
