package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete/install"
	"github.com/spf13/pflag"
)

var progBase = filepath.Base(os.Args[0])
var version = "1.0.0"

func printHelp() {
	fmt.Printf(
		`Usage: %s [flags] <expression(s)> <file(s) | stdin>

Description:
  Prints the overlapping IPs and subnets.

Examples:
  # Single expression + single file
  $ %[1]s 127.0.0.1 file1.txt

  # Single expression + multiple files
  $ %[1]s 127.0.0.1 file1.txt file2.txt file3.txt

  # Multi-expression + any files
  $ cat expression-list1.txt | %[1]s -e 127.0.0.1 -e 8.8.8.8 -e - -e expression-list2.txt file1.txt file2.txt file3.txt

Flags:
  --expression, -e
      IPs, CIDRs, and/or Ranges to be filtered. Can be used multiple times.
  --help
      Show help.

Options:
  General:
	--version, -v
	  show binary release number.
	--help, -h
	  show help.

  Completions:
	--completions-install
	  attempt completions auto-installation for any supported shell.
	--completions-bash
	  output auto-completion script for bash for manual installation.
	--completions-zsh
	  output auto-completion script for zsh for manual installation.
	--completions-fish
	  output auto-completion script for fish for manual installation.
`, progBase)
}

func cmd() error {
	var fVsn bool
	var fCompletionsInstall bool
	var fCompletionsBash bool
	var fCompletionsZsh bool
	var fCompletionsFish bool

	f := lib.CmdMatchIPFlags{}
	f.Init()
	pflag.BoolVarP(
		&fVsn,
		"version", "", false,
		"print binary release number.",
	)
	pflag.BoolVarP(
		&fCompletionsInstall,
		"completions-install", "", false,
		"attempt completions auto-installation for any supported shell.",
	)
	pflag.BoolVarP(
		&fCompletionsBash,
		"completions-bash", "", false,
		"output auto-completion script for bash for manual installation.",
	)
	pflag.BoolVarP(
		&fCompletionsZsh,
		"completions-zsh", "", false,
		"output auto-completion script for zsh for manual installation.",
	)
	pflag.BoolVarP(
		&fCompletionsFish,
		"completions-fish", "", false,
		"output auto-completion script for fish for manual installation.",
	)
	pflag.Parse()

	if fVsn {
		fmt.Println(version)
		return nil
	}

	if fCompletionsInstall {
		return install.Install(progBase)
	}
	if fCompletionsBash {
		installStr, err := install.BashCmd(progBase)
		if err != nil {
			return err
		}
		fmt.Println(installStr)
		return nil
	}
	if fCompletionsZsh {
		installStr, err := install.ZshCmd(progBase)
		if err != nil {
			return err
		}
		fmt.Println(installStr)
		return nil
	}
	if fCompletionsFish {
		installStr, err := install.FishCmd(progBase)
		if err != nil {
			return err
		}
		fmt.Println(installStr)
		return nil
	}

	return lib.CmdMatchIP(f, pflag.Args(), printHelp)
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
