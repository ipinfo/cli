package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsToolIsMulticast = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":predict.Nothing,
		"--help":  predict.Nothing,
		"-q":      predict.Nothing,
		"--quiet": predict.Nothing,
	},
}

func printHelpToolIsMultiCast(){
	fmt.Printf(`
%s tool is_multicast [<opts>] <cidr | ip | ip-range | filepath>

Description: Checks the provided address is a multicast
Inputs can be any IPv4 or IPv6 Address

Examples
  #IPv4/IPv6 Address.
  $ %[1]s tool is_multicast 127.0.0.1

  #IPv4/IPv6 Address Range
  $%[1]s tool is_multicast 127.0.0.1-127.8.95.6

  #Check CDIR
  $%[1]s tool is_multicast 127.0.0.1/32

  # Check for file.
  $ %[1]s tool is_multicast /path/to/file.txt 

  # Check entries from stdin.
  $ cat /path/to/file1.txt | %[1]s tool is_multicast

  Options:
  --help, -h
    show help.
  --quiet, -q
    quiet mode; suppress additional output.
`, progBase)
}
func cmdToolIsMultiCast() (err error){
	f := lib.CmdToolIsMulticastFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsMulticast(f, pflag.Args()[2:], printHelpToolIsMultiCast)
}