package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

func printHelpMatchIP() {
	fmt.Printf(
		`Usage: %s matchip --filter <file(s) | stdin> --criteria <file(s) | stdin>
Description:
  Prints the overlapping IPs and subnets.

Examples:
  # Match from a file
  $ %[1]s matchip --filter /path/to/list1.txt --criteria /path/to/list2.txt

  # Match from multiple files
  $ %[1]s matchip --filter=/path/to/list.txt,/path/to/list1.txt --criteria=/path/to/list2.txt,/path/to/list3.txt

  # Match from stdin
  $ cat /path/to/list1.txt | %[1]s matchip --filter - --criteria /path/to/list2.txt

Options:
  General:
    --filter, -f
      IPs, CIDRs, and/or Ranges to be filtered.
    --criteria, -c
      CIDRs and/or Ranges to check overlap with.
    --help
      show help.
`, progBase)
}

func cmdMatchIP() error {
	f := lib.CmdMatchIPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdMatchIP(f, pflag.Args()[1:], printHelpMatchIP)
}
