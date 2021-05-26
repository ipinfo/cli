package main

import (
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
)

var completions = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-o":                    predict.Nothing,
		"--only-matching":       predict.Nothing,
		"-h":                    predict.Nothing,
		"--no-filename":         predict.Nothing,
		"--no-recurse":          predict.Nothing,
		"--version":             predict.Nothing,
		"--help":                predict.Nothing,
		"--nocolor":             predict.Nothing,
		"-4":                    predict.Nothing,
		"--ipv4":                predict.Nothing,
		"-6":                    predict.Nothing,
		"--ipv6":                predict.Nothing,
		"-x":                    predict.Nothing,
		"--exclude-reserved":    predict.Nothing,
		"--completions-install": predict.Nothing,
		"--completions-bash":    predict.Nothing,
		"--completions-zsh":     predict.Nothing,
		"--completions-fish":    predict.Nothing,
	},
}

func handleCompletions() {
	completions.Complete(progBase)
}
