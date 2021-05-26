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
		`Usage: %s [<opts>] <cidr | filepath>

Description:

  Accepts CIDRs and file paths to files containing CIDRs, converting them all
  to IP ranges.

  If a file is input, it is assumed that the CIDR to convert is the first entry
  of each line (separated by '\n'). All other data remains the same.

  The IP range output is of the form "<start>-<end>" where "<start>" comes
  before or is equal to "<end>" in numeric value.

Examples:

  # Get the range for CIDR 1.1.1.0/30.
  $ %[1]s 1.1.1.0/30

  # Convert CIDR entries to IP ranges in 2 files.
  $ %[1]s /path/to/file1.txt /path/to/file2.txt

  # Convert CIDR entries to IP ranges from stdin.
  $ cat /path/to/file1.txt | %[1]s

  # Convert CIDR entries to IP ranges from stdin and a file.
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

	f := lib.CmdCIDR2RangeFlags{}
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

	return lib.CmdCIDR2Range(f, pflag.Args(), printHelp)
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
