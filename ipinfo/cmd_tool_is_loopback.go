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

func printHelpToolIsLoopBack(){
	fmt.Printf(`
%s tool is_loopback [<opts>] <cidr | ip | ip-range | filepath>

Description: Checks the provided address is a loopback
Inputs can be any IPv4 or IPv6 Address

Examples
  #IPv4/IPv6 Address.
  $ %[1]s tool is_loopback 127.0.0.1

  #IPv4/IPv6 Address Range
  $%[1]s tool is_loopback 127.0.0.1-127.8.95.6

  #Check CDIR
  $%[1]s tool is_loopback 127.0.0.1/32

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

func cmdToolIsLoopBack() (err error){
	f := lib.CmdToolIsLoopbackFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdToolIsLoopback(f, pflag.Args()[2:], printHelpToolIsLoopBack)
}
