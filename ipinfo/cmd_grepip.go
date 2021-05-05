package main

import (
	"fmt"

	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

func printHelpGrepIP() {
	fmt.Printf(
		`Usage: %s grepip [<opts>]

Options:
  General:
    --only-matching, -o
      print only matched IP in result line, excluding surrounding content.
    --no-filename, -h
      don't print source of match in result lines when more than 1 source.
    --no-recurse
      don't recurse into more directories in directory sources.
    --help
      show help.

  Outputs:
    --nocolor
      disable colored output.

  Filters:
    --ipv4, -4
      match only IPv4 addresses.
    --ipv6, -6
      match only IPv6 addresses.
    --exclude-reserved, -x
      exclude reserved/bogon IPs.
      full list can be found at https://ipinfo.io/bogon.
`, progBase)
}

func cmdGrepIP() error {
	f := lib.CmdGrepIPFlags{}
	f.Init()
	pflag.Parse()

	return lib.CmdGrepIP(f, pflag.Args()[1:], printHelpGrepIP)
}
