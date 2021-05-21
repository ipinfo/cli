package main

import (
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
)

var completions = &complete.Command{
	Sub: map[string]*complete.Command{
		"myip":       completionsMyIP,
		"bulk":       completionsBulk,
		"summarize":  completionsSummarize,
		"map":        completionsMap,
		"prips":      completionsPrips,
		"grepip":     lib.CompletionsGrepIP,
		"login":      completionsLogin,
		"logout":     completionsLogout,
		"completion": completionsCompletion,
		"version":    completionsVersion,
	},
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func handleCompletions() {
	completions.Complete(progBase)
}
