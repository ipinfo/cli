package lib

import (
	"fmt"
	"net"
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

	actionFunc := func(input string, inputType INPUT_TYPE) error {
		switch inputType {
		case INPUT_TYPE_IP:
			ActionForIPNextPrev(input, increment)
		default:
			return ErrNotIP
		}
		return nil
	}
	err := GetInputFrom(args, true, true, actionFunc)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func ActionForIPNextPrev(input string, delta int) {
	ip := net.ParseIP(input)
	if ip != nil {
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
		fmt.Println(nextPrevIP)
	}
}