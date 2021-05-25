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
		`Usage: %s [<opts>]

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
	pflag.BoolVarP(&fVsn, "version", "v", false, "print binary release number.")
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
