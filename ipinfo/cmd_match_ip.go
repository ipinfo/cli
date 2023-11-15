package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

func printHelpMatchIP() {
	fmt.Printf(
		`Usage: %s matchip [flags] <expression(s)> <file(s) | stdin>

Description:
  Prints the overlapping IPs and subnets.

Examples:
  # Single expression + single file
  $ %[1]s matchip 127.0.0.1 file1.txt
  # Single expression + multiple files
  $ %[1]s matchip 127.0.0.1 file1.txt file2.txt file3.txt
  # Multi-expression + any files
  $ cat expression-list1.txt | %[1]s matchip -e 127.0.0.1 -e 8.8.8.8 -e - -e expression-list2.txt file1.txt file2.txt file3.txt
Flags:
  --expression, -e
      IPs, CIDRs, and/or Ranges to be filtered. Can be used multiple times.
  --help
      Show help.
`, progBase)
}

func cmdMatchIP() error {
	f := lib.CmdMatchIPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdMatchIP(f, pflag.Args()[1:], printHelpMatchIP)
}
