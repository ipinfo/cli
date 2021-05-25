package main

import (
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/cli/lib/complete"
)

var completions = lib.CompletionsRange2CIDR

func handleCompletions() {
	completions.Complete(progBase)
}
