package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

var progBase = filepath.Base(os.Args[0])
var version = "1.0.0"

func printHelp() {
	fmt.Printf(
		`Usage: grepip [<opts>]

Options:
  General:
    --only-matching, -o
      print only matched IP in result line, excluding surrounding content.
    --no-filename, -h
      don't print source of match in result lines when more than 1 source.
    --no-recurse
      don't recurse into more directories in directory sources.
    --version
      show binary release number.
    --help
      show help.

  Filters:
    --ipv4, -4
      match only IPv4 addresses.
    --ipv6, -6
      match only IPv6 addresses.
    --exclude-reserved, -x
      exclude reserved/bogon IPs.
      full list can be found at https://ipinfo.io/bogon.
`)
}

func cmd() error {
	var fVsn bool

	f := lib.CmdGrepIPFlags{}
	f.Init()
	pflag.BoolVarP(&fVsn, "version", "", false, "print binary release number.")
	pflag.Parse()

	if fVsn {
		fmt.Println(version)
		return nil
	}

	return lib.CmdGrepIP(f, pflag.Args(), printHelp)
}

func main() {
	if err := cmd(); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
