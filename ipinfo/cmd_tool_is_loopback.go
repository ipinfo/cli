package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsLoopBack = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":      predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsLoopBack() {
	fmt.Printf(
		`Usage: %s tool is_loopback [<opts>] <cidr | ip | ip-range | filepath>

Description:
  Checks if the input is a loopback address.
  Inputs can be IPs, IP ranges, CIDRs, or filepath to a file

Examples:
  $ %[1]s tool is_loopback 127.0.0.0 | ::1
  $ %[1]s tool is_loopback 160.0.0.0 | fe08::2

  # Check CIDR.
  $ %[1]s tool is_loopback 127.0.0.0/32 | ::1/64
  $ %[1]s tool is_loopback 128.0.0.0/32 | fe08::2/64

  # Check IP range.
  $ %[1]s tool is_loopback 127.0.0.1-127.20.1.244
  $ %[1]s tool is_loopback 128.0.0.1-128.30.1.125

  # Check for file.
  $ %[1]s tool is_loopback /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_loopback

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolIsLoopBack() (err error) {
	f := lib.CmdToolIsLoopbackFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsLoopback(f, pflag.Args()[2:], printHelpToolIsLoopBack)
}
