package main

import (
	"github.com/ipinfo/cli/lib"
)

var completions = lib.CompletionsGrepIP

func handleCompletions() {
	completions.Complete(progBase)
}
