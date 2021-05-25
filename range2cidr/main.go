package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

var progBase = filepath.Base(os.Args[0])
var version = "1.0.0"

func printHelp() {
	fmt.Printf(
		`Usage: %s [<opts>] <ip-range | filepath>

Description:

  Accepts IP ranges and file paths to files containing IP ranges, converting
  them all to CIDRs (and multiple CIDRs if required).

  If a file is input, it is assumed that the IP range to convert is the first
  entry of each line (separated by '\n'). All other data remains the same.

  If multiple CIDRs are needed to represent an IP range on a line with other
  data, the data is copied per CIDR required. For example:

    in[0]: "1.1.1.0,1.1.1.2,other-data"
    out[0]: "1.1.1.0/31,other-data"
    out[1]: "1.1.1.2/32,other-data"

  IP ranges can of the form "<start><sep><end>" where "<sep>" can be "," or
  "-", and "<start>" and "<end>" can be any 2 IPs; order does not matter, but
  the resulting CIDRs are printed in the order they cover the range.

Examples:

  # Get all CIDRs for range 1.1.1.0-1.1.1.2.
  $ %[1]s 1.1.1.0-1.1.1.2
  $ %[1]s 1.1.1.0,1.1.1.2

  # Convert all range entries to CIDRs in 2 files.
  $ %[1]s /path/to/file1.txt /path/to/file2.txt

  # Convert all range entries to CIDRs from stdin.
  $ cat /path/to/file1.txt | %[1]s

  # Convert all range entries to CIDRs from stdin and a file.
  $ cat /path/to/file1.txt | %[1]s /path/to/file2.txt

Options:
  --version, -v
    show binary release number.
  --help, -h
    show help.
`, progBase)
}

func cmd() error {
	var fVsn bool

	f := lib.CmdRange2CIDRFlags{}
	f.Init()
	pflag.BoolVarP(
		&fVsn,
		"version", "v", false,
		"print binary release number.",
	)
	pflag.Parse()

	if fVsn {
		fmt.Println(version)
		return nil
	}

	return lib.CmdRange2CIDR(f, pflag.Args(), printHelp)
}

func main() {
	// obey NO_COLOR env var.
	if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	handleCompletions()

	if err := cmd(); err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
