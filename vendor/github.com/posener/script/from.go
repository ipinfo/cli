package script

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// From creates a stream from a reader.
func From(name string, r io.Reader) Stream {
	return Stream{stage: name, r: r}
}

// Writer creates a stream from a function that writes to a writer.
func Writer(name string, writer func(io.Writer) error) Stream {
	b := bytes.NewBuffer(nil)
	err := writer(b)
	return Stream{stage: name, r: b, err: err}
}

// Stdin starts a stream from stdin.
func Stdin() Stream {
	stdin := io.NopCloser(os.Stdin) // Prevent closing of stdin.
	return From("stdin", stdin)
}

// Echo writes to stdout.
//
// Shell command: `echo <s>`
func Echo(s string) Stream {
	return From("echo", strings.NewReader(s+"\n"))
}
