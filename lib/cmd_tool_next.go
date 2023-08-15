package lib

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

type CmdToolNextFlags struct {
	Help  bool
	Quiet bool
}

func (f *CmdToolNextFlags) Init() {
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

func CmdToolNext(
	f CmdToolNextFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	actionStdin := func(input string) {
		ActionForStdinNext(input)
	}
	actionFile := func(input string) {
		ActionForFileNext(input)
	}

	// Process inputs using the IPInputAction function.
	err := IPInputAction(args, true, true, false, false, true,
		actionStdin, nil, nil, actionFile)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionForStdinNext(input string) {
	ip := net.ParseIP(input)
	if ip != nil {
		nextIP := NextIP(ip)
		fmt.Println(nextIP)
	}
}

func ActionForFileNext(pathToFile string) {
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
		ActionForStdinNext(input) // Process each IP from the file
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func NextIP(ip net.IP) net.IP {
	nextIP := make(net.IP, len(ip))
	copy(nextIP, ip)

	for i := len(nextIP) - 1; i >= 0; i-- {
		if nextIP[i] < 255 {
			nextIP[i]++
			break
		} else {
			nextIP[i] = 0
		}
	}

	return nextIP
}
