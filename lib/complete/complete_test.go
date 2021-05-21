package complete

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCmd = &Command{
	Flags: map[string]Predictor{
		"--cmd-flag": nil,
	},
	Sub: map[string]*Command{
		"flags": {
			Flags: map[string]Predictor{
				"--values":    set{"a", "a a", "b"},
				"--something": set{""},
				"--nothing":   nil,
			},
		},
		"sub1": {
			Flags: map[string]Predictor{
				"--flag1": nil,
			},
			Sub: map[string]*Command{
				"sub11": {
					Flags: map[string]Predictor{
						"--flag11": nil,
					},
				},
				"sub12": {},
			},
			Args: set{"arg1", "arg2"},
		},
		"sub2": {},
		"args": {
			Args: set{"a", "a a", "b"},
		},
	},
}

func TestCompleter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args string
		want []string
	}{
		// Check empty flag name matching.

		{args: "flags ", want: []string{"--values", "--nothing", "--something", "--cmd-flag"}},
		{args: "flags -", want: []string{"--values", "--nothing", "--something", "--cmd-flag"}},
		{args: "flags --", want: []string{"--values", "--nothing", "--something", "--cmd-flag"}},

		// If started a flag with no matching prefix, expect nothing.
		{args: "flags -x", want: []string{}},

		// Check prefix matching for chain of sub commands.
		{args: "sub1 sub11 --f", want: []string{"--flag11", "--flag1"}},
		{args: "sub1 sub11 --fl", want: []string{"--flag11", "--flag1"}},

		// Test sub command completion.
		{args: "", want: []string{"flags", "sub1", "sub2", "args"}},
		{args: " ", want: []string{"flags", "sub1", "sub2", "args"}},
		{args: "f", want: []string{"flags"}},
		{args: "sub", want: []string{"sub1", "sub2"}},
		{args: "sub1", want: []string{"sub1"}},
		{args: "sub1 ", want: []string{"sub11", "sub12"}},

		// Suggest no sub commands if prefix is not known.
		{args: "x", want: []string{}},

		// Suggest flag value.

		// A flag that has an empty completion should return empty completion.
		// It "completes something"... But it doesn't know what, so we should
		// not complete anything else.
		{args: "flags --something ", want: []string{""}},
		{args: "flags --something foo", want: []string{""}},
		// A flag that have nil completion should complete all other options.
		{args: "flags --nothing ", want: []string{"--values", "--nothing", "--something", "--cmd-flag"}},
		// Trying to provide a value to the nothing flag should revert the phrase back to nothing.
		{args: "flags --nothing=", want: []string{}},
		// The flag value was not started, suggest all relevant values.
		{args: "flags --values ", want: []string{"a", "a\\ a", "b"}},
		{args: "flags --values a", want: []string{"a", "a\\ a"}},
		{args: "flags --values a\\", want: []string{"a\\ a"}},
		{args: "flags --values a\\ ", want: []string{"a\\ a"}},
		{args: "flags --values a\\ a", want: []string{"a\\ a"}},
		{args: "flags --values a\\ a ", want: []string{"--values", "--nothing", "--something", "--cmd-flag"}},
		{args: "flags --values \"a", want: []string{"\"a\"", "\"a a\""}},
		{args: "flags --values \"a ", want: []string{"\"a a\""}},
		{args: "flags --values \"a a", want: []string{"\"a a\""}},
		{args: "flags --values \"a a\"", want: []string{"\"a a\""}},
		{args: "flags --values \"a a\" ", want: []string{"--values", "--nothing", "--something", "--cmd-flag"}},

		{args: "flags --values=", want: []string{"a", "a\\ a", "b"}},
		{args: "flags --values=a", want: []string{"a", "a\\ a"}},
		{args: "flags --values=a\\", want: []string{"a\\ a"}},
		{args: "flags --values=a\\ ", want: []string{"a\\ a"}},
		{args: "flags --values=a\\ a", want: []string{"a\\ a"}},
		{args: "flags --values=a\\ a ", want: []string{"--values", "--nothing", "--something", "--cmd-flag"}},
		{args: "flags --values=\"a", want: []string{"\"a\"", "\"a a\""}},
		{args: "flags --values=\"a ", want: []string{"\"a a\""}},
		{args: "flags --values=\"a a", want: []string{"\"a a\""}},
		{args: "flags --values=\"a a\"", want: []string{"\"a a\""}},
		{args: "flags --values=\"a a\" ", want: []string{"--values", "--nothing", "--something", "--cmd-flag"}},

		// Complete positional arguments

		{args: "args ", want: []string{"--cmd-flag", "a", "a\\ a", "b"}},
		{args: "args a", want: []string{"a", "a\\ a"}},
		{args: "args a\\", want: []string{"a\\ a"}},
		{args: "args a\\ ", want: []string{"a\\ a"}},
		{args: "args a\\ a", want: []string{"a\\ a"}},
		{args: "args a\\ a ", want: []string{"--cmd-flag", "a", "a\\ a", "b"}},
		{args: "args \"a", want: []string{"\"a\"", "\"a a\""}},
		{args: "args \"a ", want: []string{"\"a a\""}},
		{args: "args \"a a", want: []string{"\"a a\""}},
		{args: "args \"a a\"", want: []string{"\"a a\""}},
		{args: "args \"a a\" ", want: []string{"--cmd-flag", "a", "a\\ a", "b"}},

		// Complete positional arguments from a parent command
		{args: "sub1 sub12 arg", want: []string{"arg1", "arg2"}},
	}

	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			Test(t, testCmd, tt.args, tt.want)
		})
	}
}

func TestComplete(t *testing.T) {
	defer func() {
		getEnv = os.Getenv
		exit = os.Exit
	}()

	in, out, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer func(o *os.File) { os.Stdout = o }(os.Stdout)
	defer out.Close()
	os.Stdout = out
	go io.Copy(ioutil.Discard, in)

	tests := []struct {
		line, point string
		shouldExit  bool
		shouldPanic bool
		install     string
		uninstall   string
	}{
		{shouldExit: true, line: "cmd", point: "1"},
		{shouldExit: false, line: "", point: ""},
		{shouldPanic: true, line: "cmd", point: ""},
		{shouldPanic: true, line: "cmd", point: "a"},
		{shouldExit: true, line: "cmd", point: "4"},

		{shouldExit: true, install: "1"},
		{shouldExit: false, install: "a"},
		{shouldExit: true, uninstall: "1"},
		{shouldExit: false, uninstall: "a"},
	}

	for _, tt := range tests {
		t.Run(tt.line+"@"+tt.point, func(t *testing.T) {
			getEnv = func(env string) string {
				switch env {
				case "COMP_LINE":
					return tt.line
				case "COMP_POINT":
					return tt.point
				case "COMP_INSTALL":
					return tt.install
				case "COMP_UNINSTALL":
					return tt.uninstall
				case "COMP_YES":
					return "0"
				default:
					panic(env)
				}
			}
			isExit := false
			exit = func(int) {
				isExit = true
			}
			if tt.shouldPanic {
				assert.Panics(t, func() { testCmd.Complete("") })
			} else {
				testCmd.Complete("")
				assert.Equal(t, tt.shouldExit, isExit)
			}
		})
	}
}

// ExampleComplete_outputCapturing demonstrates the ability to capture
// the output of Complete() invocations, crucial for integration tests.
func ExampleComplete_outputCapturing() {
	defer func(f func(int)) { exit = f }(exit)
	defer func(f getEnvFn) { getEnv = f }(getEnv)
	exit = func(int) {}

	// This is where the actual example starts:

	cmd := &Command{Sub: map[string]*Command{"bar": {}}}
	getEnv = promptEnv("foo b")

	Complete("foo", cmd)

	// Output:
	// bar
}

type set []string

func (s set) Predict(_ string) []string {
	return s
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		s      string
		prefix string
		want   string
		wantOK bool
	}{
		{s: "ab", prefix: `b`, want: ``, wantOK: false},
		{s: "", prefix: `b`, want: ``, wantOK: false},
		{s: "ab", prefix: `a`, want: `ab`, wantOK: true},
		{s: "ab", prefix: `"'b`, want: ``, wantOK: false},
		{s: "ab", prefix: `"'a`, want: `"'ab'"`, wantOK: true},
		{s: "ab", prefix: `'"a`, want: `'"ab"'`, wantOK: true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"/"+tt.prefix, func(t *testing.T) {
			got, gotOK := hasPrefix(tt.s, tt.prefix)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOK, gotOK)
		})
	}
}

// getEnvFn emulates os.GetEnv by mapping one string to another.
type getEnvFn = func(string) string

// promptEnv returns getEnvFn that emulates the environment variables
// a shell would set when its prompt has the given contents.
var promptEnv = func(contents string) getEnvFn {
	return func(key string) string {
		switch key {
		case "COMP_LINE":
			return contents
		case "COMP_POINT":
			return strconv.Itoa(len(contents))
		}
		return ""
	}
}
