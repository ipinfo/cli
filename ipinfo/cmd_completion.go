package main

import (
	"fmt"
	"strings"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/install"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsCompletion = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
	Args: predict.Set([]string{
		"install",
		"bash",
		"zsh",
		"fish",
	}),
}

func printHelpCompletion() {
	fmt.Printf(
		`Usage: %s completion [<opts>] [install | bash | zsh | fish]

Description:
  Install or print out the code needed to do shell auto-completion.

  The current explicitly supported shells are:
  - bash
  - zsh
  - fish

Examples:
  # Attempt auto-installation on any of the supported shells.
  $ %[1]s completion install

  # Output auto-completion script for bash for manual installation.
  $ %[1]s completion bash

  # Output auto-completion script for zsh for manual installation.
  $ %[1]s completion zsh

  # Output auto-completion script for fish for manual installation.
  $ %[1]s completion fish

Options:
  --help, -h
    show help.
`, progBase)
}

func cmdCompletion() error {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	args := pflag.Args()[1:]
	if fHelp || len(args) != 1 {
		printHelpCompletion()
		return nil
	}

	var installStr string
	var err error
	switch strings.ToLower(args[0]) {
	case "install":
		return install.Install(progBase)
	case "bash":
		installStr, err = install.BashCmd(progBase)
	case "zsh":
		installStr, err = install.ZshCmd(progBase)
	case "fish":
		installStr, err = install.FishCmd(progBase)
	default:
		fmt.Printf("err: %s is not a valid subcommand\n\n", args[0])
		printHelpCompletion()
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Println(installStr)

	return nil
}
