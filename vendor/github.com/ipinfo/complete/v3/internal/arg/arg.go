package arg

import (
	"strings"

	"github.com/ipinfo/complete/v3/internal/tokener"
)

// Arg is typed a command line argument.
type Arg struct {
	Text      string
	Completed bool
	Parsed
}

// Parsed contains information about the argument.
type Parsed struct {
	Flag     string
	HasFlag  bool
	Value    string
	HasValue bool
}

// Parse parses a typed command line argument list, and returns a list of
// arguments.
func Parse(line string) []Arg {
	var args []Arg
	for {
		arg, after := next(line)
		if arg.Text != "" {
			args = append(args, arg)
		}
		line = after
		if line == "" {
			break
		}
	}
	return args
}

// next returns the first argument in the line and the rest of the line.
func next(line string) (arg Arg, after string) {
	defer arg.parse()
	// Start and end of the argument term.
	var start, end int

	// Stack of quote marks met during the paring of the argument.
	var token tokener.Tokener

	// Skip prefix spaces.
	for start = 0; start < len(line); start++ {
		token.Visit(line[start])
		if !token.LastSpace() {
			break
		}
	}

	// If line is only spaces, return empty argument and empty leftovers.
	if start == len(line) {
		return
	}

	for end = start + 1; end < len(line); end++ {
		token.Visit(line[end])
		if token.LastSpace() {
			arg.Completed = true
			break
		}
	}
	arg.Text = line[start:end]
	if !arg.Completed {
		return
	}
	start2 := end

	// Skip space after word.
	for start2 < len(line) {
		token.Visit(line[start2])
		if !token.LastSpace() {
			break
		}
		start2++
	}
	after = line[start2:]
	return
}

// parse a flag from an argument. The flag can have value attached when it is
// given in the `--key=value` format.
func (a *Arg) parse() {
	if len(a.Text) == 0 {
		return
	}

	// no flag; just value.
	if a.Text[0] != '-' {
		a.Value = a.Text
		a.HasValue = true
		return
	}

	// now we're sure there's at least a flag.
	a.HasFlag = true

	// get '=' sign to see if we have a value.
	eqIdx := strings.Index(a.Text, "=")
	if eqIdx == -1 {
		a.HasValue = false
	} else if eqIdx != len(a.Text)-1 {
		a.HasValue = true
		a.Value = a.Text[eqIdx+1:]
	} else {
		a.HasValue = true
	}

	// pull out flag value if any yet.
	if eqIdx == -1 {
		eqIdx = len(a.Text)
	}
	a.Flag = a.Text[:eqIdx]

	// no flag name isn't valid, e.g. '--='.
	if a.HasValue && a.Flag[len(a.Flag)-1] == '-' {
		a.Parsed = Parsed{}
		return
	}
}
