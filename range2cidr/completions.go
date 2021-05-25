package main

import (
	"github.com/ipinfo/cli/lib"
)

var completions = lib.CompletionsRange2CIDR

func handleCompletions() {
	completions.Complete(progBase)
}
