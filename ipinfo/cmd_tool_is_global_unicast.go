package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsGlobalUnicast = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":predict.Nothing,
		"--help":predict.Nothing,
		"-q":predict.Nothing,
		"--quiet":predict.Nothing,
	},
}

func printHelpToolIsGlobalUnicast(){
	fmt.Printf(
		`Usage: %s tool is_global_unicast [<opts>] <cidr | ip | ip-range | filepath>

Description:
  Checks if the input is an globl unicast address.
  Inputs can be IPs, IP ranges, CIDRs, or filepath to a file

Examples:
  # Check CIDR.
  $ %[1]s tool is_global_unicast 1.1.1.0/30

  # Check IP range.
  $ %[1]s tool is_global_unicast 1.1.1.0-1.1.1.244

  # Check for file.
  $ %[1]s tool is_global_unicast /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_global_unicast

Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}

func cmdToolisGlobalUnicast()(err error){
	f:=lib.CmdToolIsGlobalUnicastFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsGlobalUnicast(f,pflag.Args()[2:],printHelpToolIsGlobalUnicast)
}
