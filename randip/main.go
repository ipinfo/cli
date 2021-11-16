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
		`Usage: %s [<opts>] 

Description:
  Generates random IP/IPs.
  By default, generates 1 random IPv4 address with starting range 0.0.0.0 and 
  ending range 255.255.255.255, but can be configured to generate any number of 
  a combination of IPv4/IPv6 addresses within any range.
  
  Using --ipv6 or -6 without any starting or ending range will generate a IP 
  between range of :: to ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff.

  Note that only IPv4 or IPv6 IPs can be generated, but not both.

  $ %[1]s  --ipv6 --count 5
  $ %[1]s  -4 -n 10
  $ %[1]s  -4 -s 1.1.1.1 -e 10.10.10.10
  $ %[1]s  -6 --start 9c61:f71e:656d:d12e:98a3:9814:38cf:5592
  $ %[1]s  -6 --end eedd:8977:56d9:aac3:947b:29cc:78ea:deab

Options:
  --help, -h
    show help.
  --count, -n 
    number of IPs to generate.
  --ipv4, -4
    generates IPv4 IPs.
  --ipv6, -6
    generates IPv6 IPs.
  --start, -s 
    starting range of IPs.
	default: minimum IP possible for IP type selected.
  --end, -e
    ending range of IPs.
	default: maximum IP possible for IP type selected.

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

	f := lib.CmdRandIPFlags{}
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

	return lib.CmdRandIP(f, pflag.Args(), printHelp)
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
