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

	increment := 1

	actionStdin := func(input string) {
		ActionForStdinNextPrev(input, increment)
	}
	actionFile := func(input string) {
		ActionForFileNextPrev(input, increment)
	}

	err := IPInputAction(args, true, true, false, false, true,
		actionStdin, nil, nil, actionFile)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func ActionForStdinNextPrev(input string, delta int) {
	ip := net.ParseIP(input)
	if ip != nil {
		nextPrevIP := UpdateIPAddress(ip, delta)
		fmt.Println(nextPrevIP)
	}
}

func ActionForFileNextPrev(pathToFile string, delta int) {
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
		ActionForStdinNextPrev(input, delta) // Process each IP from the file
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
}

func UpdateIPAddress(ip net.IP, delta int) net.IP {
	nextPrevIP := make(net.IP, len(ip))
	copy(nextPrevIP, ip)

	for i := len(nextPrevIP) - 1; i >= 0; i-- {
		if delta > 0 {
			if nextPrevIP[i] < 255-byte(delta) {
				nextPrevIP[i] += byte(delta)
				break
			} else {
				nextPrevIP[i] = 255
			}
		} else if delta < 0 {
			if nextPrevIP[i] >= byte(-delta) {
				nextPrevIP[i] += byte(delta)
				break
			} else {
				nextPrevIP[i] = 0
			}
		}
	}
	return nextPrevIP
}
