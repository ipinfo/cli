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
		`Usage: %s [<opts>] <ip | ip-range | cidr | file>

Description:
  Accepts CIDRs (e.g. 8.8.8.0/24) and IP ranges (e.g. 8.8.8.0-8.8.8.255).

Examples:
  # List all IPs in a CIDR.
  $ %[1]s 8.8.8.0/24

  # List all IPs in multiple CIDRs.
  $ %[1]s 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # List all IPs in an IP range.
  $ %[1]s 8.8.8.0-8.8.8.255

  # List all IPs in multiple CIDRs and IP ranges.
  $ %[1]s 1.1.1.0/30 8.8.8.0-8.8.8.255 2.2.2.0/30 7.7.7.0,7.7.7.10

  # List all IPs from stdin input (newline-separated).
  $ echo -e '1.1.1.0/30\n8.8.8.0-8.8.8.255\n7.7.7.0,7.7.7.10' | %[1]s

Options:
  --help, -h
    show help.
`, progBase)
}

func cmd() error {
	var fVsn bool
	var fCompletionsInstall bool
	var fCompletionsBash bool
	var fCompletionsZsh bool
	var fCompletionsFish bool

	f := lib.CmdPripsFlags{}
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

	return lib.CmdPrips(f, pflag.Args(), printHelp)
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
