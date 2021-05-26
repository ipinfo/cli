package script

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/hashicorp/go-multierror"
)

// Exec executes a command and returns a stream of the stdout of the command.
func Exec(cmd string, args ...string) Stream {
	return From("empty", nil).Through(exe{cmd: cmd, args: args})
}

// ExecHandleStderr executes a command, returns a stream of the stdout of the command and enable
// collecting the stderr of the command.
//
// If the stderr is nil, it will be ignored.
//
// For example, collecting the stderr to memory can be done by providing a `&bytes.Buffer` as
// `stderr`. Writing it to stderr can be done by providing `os.Stderr` as `stderr`. Logging it
// to a file can be done by providing an `os.File` as the `stderr`.
func ExecHandleStderr(stderr io.Writer, cmd string, args ...string) Stream {
	return From("empty", nil).Through(exe{cmd: cmd, args: args, stderr: stderr})
}

// Exec executes a command and returns a stream of the stdout of the command.
func (s Stream) Exec(cmd string, args ...string) Stream {
	return s.Through(exe{cmd: cmd, args: args})
}

// ExecHandleStderr executes a command, returns a stream of the stdout of the command and enable
// collecting the stderr of the command.
//
// If the stderr is nil, it will be ignored.
func (s Stream) ExecHandleStderr(stderr io.Writer, cmd string, args ...string) Stream {
	return s.Through(exe{cmd: cmd, args: args, stderr: stderr})
}

type exe struct {
	cmd    string
	args   []string
	stderr io.Writer
}

func (e exe) Name() string {
	return fmt.Sprintf("exec(%v, %+v)", e.cmd, e.args)
}

func (e exe) Pipe(stdin io.Reader) (io.Reader, error) {
	cmd := exec.Command(e.cmd, e.args...)
	var errors *multierror.Error

	// Pipe previous stdin if available.
	if stdin != nil {
		cmd.Stdin = stdin
	}

	// Pipe stdout to the current command output.
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		errors = multierror.Append(errors, fmt.Errorf("pipe stdout: %v", err))
	}

	if e.stderr == nil {
		e.stderr = ioutil.Discard
	}
	cmd.Stderr = e.stderr

	// start the process
	err = cmd.Start()
	if err != nil {
		errors = multierror.Append(errors, fmt.Errorf("start process: %v", err))
	}
	return readcloser{
		Reader: cmdOut,
		Closer: closerFn(func() error { return cmd.Wait() }),
	}, errors.ErrorOrNil()
}

type closerFn func() error

func (f closerFn) Close() error { return f() }
