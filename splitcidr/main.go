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
var version = "1.2.0"

func printHelp() {
	fmt.Printf(
		`Usage: %s splitcidr <cidr> <split>

Description:
  splits a larger CIDR into smaller CIDRs.

  $ %[1]s splitcidr 8.8.8.0/24 25


Options:
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

	f := lib.CmdSplitCIDRFlags{}
	f.Init()
	pflag.BoolVarP(
		&fVsn,
		"version", "v", false,
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

	return lib.CmdSplitCIDR(f, pflag.Args(), printHelp)
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
