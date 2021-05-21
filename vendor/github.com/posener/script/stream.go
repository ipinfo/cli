// Package script provides helper functions to write scripts.
//
// Inspired by https://github.com/bitfield/script, with some improvements:
//
// * Output between streamed commands is a stream and not loaded to memory.
//
// * Better representation and handling of errors.
//
// * Proper incocation, usage and handling of stderr of custom commands.
//
// The script chain is represented by a
// (`Stream`) https://godoc.org/github.com/posener/script#Stream object. While each command in the
// stream is abstracted by the (`Command`) https://godoc.org/github.com/posener/script#Command
// struct. This library provides basic functionality, but can be extended freely.
package script

import (
	"io"
	"reflect"

	"github.com/hashicorp/go-multierror"
)

// Stream is a chain of operations on a stream of bytes. The stdout of each operation in the stream
// feeds the following operation stdin. The stream object have different methods that allow
// manipulating it, most of them resemble well known linux commands.
//
// A custom modifier can be used with the `Through` or with the `Modify` functions.
//
// The stream object is created by some in this library or from any `io.Reader` using the `From`
// function. It can be dumped using some functions in this library, to a custom reader using the
// `To` method.
type Stream struct {
	// r is the output reader of this stream. If r also implements the `io.Closer` interface, it
	// will be closed when the stream is closed.
	r io.Reader
	// stage is the name of the current stage in the stream.
	stage string
	// parent points to the stage before the current stage in the stream.
	parent *Stream
	// err contains an error from the current stage in the stream.
	err error
}

// Read can be used to read from the stream.
func (s Stream) Read(b []byte) (int, error) {
	return s.r.Read(b)
}

// Close closes all the stages in the stream and return the errors that occurred in all of the
// stages.
func (s Stream) Close() error {
	var errors *multierror.Error
	for cur := &s; cur != nil; cur = cur.parent {
		if cur.err != nil {
			errors = multierror.Append(errors, cur.err)
		}
		if closer, ok := cur.r.(io.Closer); ok {
			if err := closer.Close(); err != nil {
				errors = multierror.Append(errors, err)
			}
		}
	}
	return errors.ErrorOrNil()
}

// Through passes the current stream through a pipe. This function can be used to add custom
// commands that are not available in this library.
func (s Stream) Through(pipe Pipe) Stream {
	r, err := pipe.Pipe(s.r)
	if r == nil {
		panic("a command must contain a reader")
	}
	return Stream{
		stage:  pipe.Name(),
		r:      r,
		err:    err,
		parent: &s,
	}
}

// Pipe reads from a reader and returns another reader.
type Pipe interface {
	// Pipe gets a reader and returns another reader. A pipe may return an error and a reader
	// together.
	Pipe(stdin io.Reader) (io.Reader, error)
	// Name of pipe.
	Name() string
}

// PipeFn is a function that implements Pipe.
type PipeFn func(io.Reader) (io.Reader, error)

func (f PipeFn) Pipe(stdin io.Reader) (io.Reader, error) { return f(stdin) }

func (f PipeFn) Name() string { return reflect.TypeOf(f).Name() }

type readcloser struct {
	io.Reader
	io.Closer
}
