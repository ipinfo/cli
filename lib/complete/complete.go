package complete

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ipinfo/cli/lib/complete/install"
)

// Completer is an interface that a command line should implement in order to
// get bash completion.
type Completer interface {
	// SubCmdList should return the list of all sub commands of the current
	// command.
	SubCmdList() []string

	// SubCmdGet should return a sub command of the current command for the
	// given sub command name.
	SubCmdGet(cmd string) Completer

	// FlagList should return a list of all the flag names of the current
	// command. The flag names should not have the dash prefix.
	FlagList() []string

	// FlagGet should return completion options for a given flag. It is invoked
	// with the flag name without the dash prefix. The flag is not promised to
	// be in the command flags. In that case, this method should return a nil
	// predictor.
	FlagGet(flag string) Predictor

	// ArgsGet should return predictor for positional arguments of the command
	// line.
	ArgsGet() Predictor
}

// Predictor can predict completion options.
type Predictor interface {
	// Predict returns prediction options for a given prefix. The prefix is
	// what currently is typed as a hint for what to return, but the returned
	// values can have any prefix. The returned values will be filtered by the
	// prefix when needed regardless. The prefix may be empty which means that
	// no value was typed.
	Predict(prefix string) []string
}

// PredictFunc is a function that implements the Predictor interface.
type PredictFunc func(prefix string) []string

func (p PredictFunc) Predict(prefix string) []string {
	if p == nil {
		return nil
	}
	return p(prefix)
}

var (
	getEnv = os.Getenv
	exit   = os.Exit
)

// Complete the command line arguments for the given command in the case that
// the program was invoked with COMP_LINE and COMP_POINT environment variables.
// In that case it will also `os.Exit()`. The program name should be provided
// for installation purposes.
func Complete(name string, cmd Completer) {
	var (
		line        = getEnv("COMP_LINE")
		point       = getEnv("COMP_POINT")
		doInstall   = getEnv("COMP_INSTALL") == "1"
		doUninstall = getEnv("COMP_UNINSTALL") == "1"
		yes         = getEnv("COMP_YES") == "1"
	)
	var (
		out io.Writer = os.Stdout
		in  io.Reader = os.Stdin
	)
	if doInstall || doUninstall {
		install.Run(name, doUninstall, yes, out, in)
		exit(0)
		return
	}
	if line == "" {
		return
	}
	i, err := strconv.Atoi(point)
	if err != nil {
		panic("COMP_POINT env should be integer, got: " + point)
	}
	if i > len(line) {
		i = len(line)
	}

	// Parse the command line up to the completion point.
	args := Parse(line[:i])

	// The first word is the current command name.
	args = args[1:]

	// Run the completion algorithm.
	options, err := completer{Completer: cmd, args: args}.complete()
	if err != nil {
		fmt.Fprintln(out, "\n"+err.Error())
	} else {
		for _, option := range options {
			fmt.Fprintln(out, option)
		}
	}
	exit(0)
}

type completer struct {
	Completer
	args  []Arg
	stack []Completer
}

// complete command with given before and after text.
// if the command has sub commands: try to complete only sub commands.
// Otherwise complete flags and positional arguments.
func (c completer) complete() ([]string, error) {
reset:
	arg := Arg{}
	if len(c.args) > 0 {
		arg = c.args[0]
	}

	switch {
	case len(c.SubCmdList()) == 0:
		// No sub commands, parse flags and positional arguments.
		return c.suggestLeafCommandOptions(), nil
	case !arg.Completed:
		// Currently typing a sub command.
		return c.suggestSubCommands(arg.Text), nil
	case c.SubCmdGet(arg.Text) != nil:
		// Sub command completed, look into that sub command completion.
		// Set the complete command to the requested sub command, and the
		// before text to all the text after the command name and rerun the
		// complete algorithm with the new sub command.
		c.stack = append([]Completer{c.Completer}, c.stack...)
		c.Completer = c.SubCmdGet(arg.Text)
		c.args = c.args[1:]
		goto reset
	default:
		return nil, nil
	}
}

func (c completer) suggestSubCommands(prefix string) []string {
	subs := c.SubCmdList()
	return suggest(prefix, func(prefix string) []string {
		var options []string
		for _, sub := range subs {
			if strings.HasPrefix(sub, prefix) {
				options = append(options, sub)
			}
		}
		return options
	})
}

func (c completer) suggestLeafCommandOptions() (options []string) {
	arg, before := Arg{}, Arg{}
	if len(c.args) > 0 {
		arg = c.args[len(c.args)-1]
	}
	if len(c.args) > 1 {
		before = c.args[len(c.args)-2]
	}

	if arg.Completed {
		// on a flag such as `--flag <val>` where `val` is currently empty?
		if arg.HasFlag && !arg.HasValue {
			if opts := c.suggestFlagValue(arg.Flag, arg.Value); opts != nil {
				return opts
			}
		}

		return []string{}
	}

	// Complete value being typed.
	if arg.HasValue {
		// Complete value of current flag.
		if arg.HasFlag {
			return c.suggestFlagValue(arg.Flag, arg.Value)
		}
		// Complete value of flag in a previous argument.
		if before.HasFlag && !before.HasValue {
			return c.suggestFlagValue(before.Flag, arg.Value)
		}
	}

	// no val; suggest a flag if dashes started.
	if !arg.HasValue && arg.HasFlag {
		options = c.suggestFlag(arg.Flag)
	}

	// no flag; suggest positional arg.
	if !arg.HasFlag {
		options = append(options, c.suggestArgsValue(arg.Value)...)
	}

	return options
}

func (c completer) suggestFlag(prefix string) []string {
	return suggest(prefix, func(prefix string) []string {
		var options []string
		c.iterateStack(func(cmd Completer) {
			// Suggest all flags with the given prefix.
			for _, name := range cmd.FlagList() {
				if strings.HasPrefix(name, prefix) {
					options = append(options, name)
				}
			}
		})
		return options
	})
}

func (c completer) suggestFlagValue(flagName, prefix string) []string {
	var options []string
	c.iterateStack(func(cmd Completer) {
		if len(options) == 0 {
			if p := cmd.FlagGet(flagName); p != nil {
				options = p.Predict(prefix)
			}
		}
	})
	return filterByPrefix(prefix, options...)
}

func (c completer) suggestArgsValue(prefix string) []string {
	var options []string
	c.iterateStack(func(cmd Completer) {
		if len(options) == 0 {
			if p := cmd.ArgsGet(); p != nil {
				options = p.Predict(prefix)
			}
		}
	})
	return filterByPrefix(prefix, options...)
}

func (c completer) iterateStack(f func(Completer)) {
	for _, cmd := range append([]Completer{c.Completer}, c.stack...) {
		f(cmd)
	}
}

func suggest(
	prefix string,
	collect func(prefix string) []string,
) []string {
	options := collect(prefix)
	if len(options) > 0 {
		return options
	}

	// Nothing matched.
	options = collect("")
	return options
}

func filterByPrefix(prefix string, options ...string) []string {
	var filtered []string
	for _, option := range options {
		if fixed, ok := hasPrefix(option, prefix); ok {
			filtered = append(filtered, fixed)
		}
	}
	if len(filtered) > 0 {
		return filtered
	}
	return options
}

// hasPrefix checks if s has the give prefix. It disregards quotes and escaped
// spaces, and return s in the form of the given prefix.
func hasPrefix(s, prefix string) (string, bool) {
	var (
		token  Tokener
		si, pi int
	)
	for ; pi < len(prefix); pi++ {
		token.Visit(prefix[pi])
		lastQuote := !token.Escaped() &&
			(prefix[pi] == '"' || prefix[pi] == '\'')

		if lastQuote {
			continue
		}
		if si == len(s) {
			break
		}
		if s[si] == ' ' && !token.Quoted() && token.Escaped() {
			s = s[:si] + "\\" + s[si:]
		}
		if s[si] != prefix[pi] {
			return "", false
		}
		si++
	}

	if pi < len(prefix) {
		return "", false
	}

	for ; si < len(s); si++ {
		token.Visit(s[si])
	}

	return token.Closed(), true
}
