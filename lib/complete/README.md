# complete

[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/ipinfo/cli/lib/complete)

Writing bash completion scripts is a hard work, usually done in the bash
scripting language. This package can help produce a binary which will serve as
the auto-completion input/output facility: your shell will run your binary
everytime completion is required with the input being the user's current
command line, and the binary's output being the suggested completions.

## Installation

Supported shells:

- [x] bash
- [x] zsh
- [x] fish

The installation of completion for a command line tool is done automatically by
this library by running the command line tool with the `COMP_INSTALL`
environment variable set. Uninstalling the completion is similarly done by the
`COMP_UNINSTALL` environment variable.  For example, if a tool called `my-cli`
uses this library, the completion can install by running
`COMP_INSTALL=1 my-cli`.

## Example

```go
import (
 	"flag"
 	"github.com/ipinfo/cli/lib/complete"
 	"github.com/ipinfo/cli/lib/complete/predict"
)

var (
 	// Add variables to the program.
 	name      = flag.String("name", "", "")
 	something = flag.String("something", "", "")
 	nothing   = flag.String("nothing", "", "")
)

func main() {
 	// Create the complete command.
 	// Here we define completion values for each flag.
 	cmd := &complete.Command{
	 	Flags: map[string]complete.Predictor{
 			"name":      predict.Set{"foo", "bar", "foo bar"},
 			"something": predict.Something,
 			"nothing":   predict.Nothing,
 		},
 	}
 	// Run the completion - provide it with the binary name.
 	cmd.Complete("my-program")
 	// Parse the flags.
 	flag.Parse()
 	// Program logic...
}
```
