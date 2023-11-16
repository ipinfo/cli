package main

import (
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
)

var completions = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-e":           predict.Nothing,
		"--expression": predict.Nothing,
		"--help":       predict.Nothing,
	},
}

func handleCompletions() {
	completions.Complete(progBase)
}
