package main

import (
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
		"grepip":     completionsGrepIP,
		"cidr2range": completionsCIDR2Range,
		"range2cidr": completionsRange2CIDR,
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
