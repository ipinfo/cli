package main

import (
	"github.com/ipinfo/cli/lib"
)

var completions = lib.CompletionsCIDR2Range

func handleCompletions() {
	completions.Complete(progBase)
}
