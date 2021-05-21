package script

import (
	"bufio"
	"io"
	"reflect"
)

// Modifier modifies input lines to output. On each line of the input the Modify method is called,
// and the modifier can change it, omit it, or break the iteration.
type Modifier interface {
	// Modify a line. The input of this function will always be a single line from the input of the
	// stream, without the trailing '\n'. It should return the output of the stream and should
	// append a trailing '\n' if it want it to be a line in the output.
	//
	// When EOF of input stream is met, the function will be called once more with a nil line value
	// to enable output any buffered data.
	//
	// When the return modified value is nil, the line will be discarded.
	//
	// When the returned eof value is true, the Read will return that error.
	Modify(line []byte) (modifed []byte, err error)
	// Name returns the name of the command that will represent this modifier.
	Name() string
}

// ModifyFn is a function for modifying input lines.
type ModifyFn func(line []byte) (modifed []byte, err error)

func (m ModifyFn) Modify(line []byte) (modifed []byte, err error) { return m(line) }

func (m ModifyFn) Name() string { return reflect.TypeOf(m).Name() }

// Modify applies modifier on every line of the input.
func (s Stream) Modify(modifier Modifier) Stream {
	return s.Through(modPipe{Modifier: modifier})
}

// modPipe takes a Modifier and exposes the Pipe interface.
type modPipe struct {
	Modifier
	r *bufio.Reader
	// partialOut stores leftover of a line that was not fully read by output.
	partialOut []byte
	err        error
}

func (m modPipe) Pipe(stdin io.Reader) (io.Reader, error) {
	m.r = bufio.NewReader(stdin)
	return &m, nil
}

func (m modPipe) Close() error {
	if m.err == io.EOF {
		return nil
	}
	return m.err
}

func (m *modPipe) Read(out []byte) (n int, err error) {
	if len(m.partialOut) > 0 {
		m.partialOut, n = copyBytes(out, m.partialOut)
		return n, nil
	}
	if m.err != nil {
		return 0, m.err
	}

	// partialIn stores a line that was not fully read from input.
	var partialIn []byte

	for {
		line, isPrefix, err := m.r.ReadLine()
		if err != nil {
			if err != io.EOF {
				return 0, err
			}
			// Remember that we have EOF for next read call.
			m.err = io.EOF
		}
		if len(partialIn) > 0 {
			line = append(partialIn, line...)
			partialIn = nil
		}
		if isPrefix {
			partialIn = line
			continue
		}

		line, err = m.Modifier.Modify(line)
		if err != nil {
			m.err = err
		}

		m.partialOut, n = copyBytes(out, line)
		return n, nil
	}
}

func copyBytes(dst, src []byte) (leftover []byte, n int) {
	n = len(src)
	if n > len(dst) {
		n = len(dst)
	}
	copy(dst[:n], src[:n])
	return src[n:], n
}
