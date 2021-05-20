# complete

[![Build Status](https://travis-ci.org/ipinfo/complete/v3.svg?branch=master)](https://travis-ci.org/ipinfo/complete/v3)
[![codecov](https://codecov.io/gh/ipinfo/complete/v3/branch/master/graph/badge.svg)](https://codecov.io/gh/ipinfo/complete/v3)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/ipinfo/complete/v3)

Package complete is everything for bash completion and Go.

Writing bash completion scripts is a hard work, usually done in the bash scripting language.
This package provides:

* A library for bash completion for Go programs.
* A tool for writing bash completion script in the Go language. For any Go or non Go program.
* Enables an easy way to install/uninstall the completion of the command.

The library and tools are extensible such that any program can add its one logic, completion types
or methologies.

## Complete Package

Supported shells:

- [x] bash
- [x] zsh
- [x] fish

The installation of completion for a command line tool is done automatically by this library by
running the command line tool with the `COMP_INSTALL` environment variable set. Uninstalling the
completion is similarly done by the `COMP_UNINSTALL` environment variable.
For example, if a tool called `my-cli` uses this library, the completion can install by running
`COMP_INSTALL=1 my-cli`.

## Usage

Add bash completion capabilities to any Go program. See [./example/command](./example/command).

```go
 import (
 	"flag"
 	"github.com/ipinfo/complete/v3"
 	"github.com/ipinfo/complete/v3/predict"
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

This package also enables to complete flags defined by the standard library `flag` package.
To use this feature, simply call `complete.CommandLine` before `flag.Parse`. (See [./example/stdlib](./example/stdlib)).

```diff
 import (
 	"flag"
+	"github.com/ipinfo/complete/v3"
 )
 var (
 	// Define flags here...
 	foo = flag.Bool("foo", false, "")
 )
 func main() {
 	// Call command line completion before parsing the flags - provide it with the binary name.
+	complete.CommandLine("my-program")
 	flag.Parse()
 }
```

If flag value completion is desired, it can be done by providing the standard library `flag.Var`
function a `flag.Value` that also implements the `complete.Predictor` interface.

## Testing

For command line bash completion testing use the `complete.Test` function.

## Sub Packages

* [example/command](./example/command): command shows how to have bash completion to an arbitrary Go program using the `complete.Command` struct.
* [example/stdlib](./example/stdlib): stdlib shows how to have flags bash completion to an arbitrary Go program that uses the standard library flag package.
* [install](./install): Package install provide installation functions of command completion.
* [predict](./predict): Package predict provides helper functions for completion predictors.

## Examples

### OutputCapturing

ExampleComplete_outputCapturing demonstrates the ability to capture
the output of Complete() invocations, crucial for integration tests.

```golang
defer func(f func(int)) { exit = f }(exit)
defer func(f getEnvFn) { getEnv = f }(getEnv)
exit = func(int) {}

// This is where the actual example starts:

cmd := &Command{Sub: map[string]*Command{"bar": {}}}
getEnv = promptEnv("foo b")

Complete("foo", cmd)
```

 Output:

```
bar
```

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
