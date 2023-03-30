package main

import (
	"github.com/ipinfo/cli/lib/complete"
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
		"-v":         nil,
		"--version":  nil,
		"-h":         nil,
		"--help":     nil,
	},
}

func handleCompletions() {
	completions.Complete(progBase)
}
