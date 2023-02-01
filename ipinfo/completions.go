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
		"cidr2ip":    completionsCIDR2IP,
		"range2cidr": completionsRange2CIDR,
		"range2ip":   completionsRange2IP,
		"randip":     completionsRandIP,
		"splitcidr":  completionsSplitCIDR,
		"cache":      completionsCache,
		"quota":      completionsQuota,
		"login":      completionsLogin,
		"logout":     completionsLogout,
		"config":     completionsConfig,
		"completion": completionsCompletion,
		"version":    completionsVersion,
	},
	Flags: map[string]complete.Predictor{
		"-v":        predict.Nothing,
		"--version": predict.Nothing,
		"-h":        predict.Nothing,
		"--help":    predict.Nothing,
	},
}

func handleCompletions() {
	completions.Complete(progBase)
}
