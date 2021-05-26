package main

import (
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
)

var completions = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-v":                    predict.Nothing,
		"--version":             predict.Nothing,
		"-h":                    predict.Nothing,
		"--help":                predict.Nothing,
		"--completions-install": predict.Nothing,
		"--completions-bash":    predict.Nothing,
		"--completions-zsh":     predict.Nothing,
		"--completions-fish":    predict.Nothing,
	},
}

func handleCompletions() {
	completions.Complete(progBase)
}
