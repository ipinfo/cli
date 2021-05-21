package script

import (
	"bytes"
	"fmt"
	"io"
)

// Head reads only the n first lines of the given reader. If n is a negative number, all lines
// besides last n lines will be read.
//
// Shell command: `head -n <n>`
func (s Stream) Head(n int) Stream {
	if n < 0 {
		return s.Modify(&negHead{n: -n, lines: make([][]byte, 0, -n)})
	}
	return s.Modify(&head{n: n})
}

// Tail reads only the n last lines of the given reader. If n is a negative number, all lines
// besides the first n lines will be read.
//
// Shell command: `tail -n <n>`
func (s Stream) Tail(n int) Stream {
	if n < 0 {
		return s.Modify(&negTail{n: -n})
	}
	return s.Modify(&tail{n: n, lines: make([][]byte, 0, n)})
}

type head struct {
	n int
}

func (h *head) Modify(line []byte) ([]byte, error) {
	if line == nil || h.n == 0 {
		return nil, io.EOF
	}
	h.n--
	return append(line, '\n'), nil
}

func (h *head) Name() string {
	return fmt.Sprintf("head(%d)", h.n)
}

type negHead struct {
	n     int
	lines [][]byte
}

func (h *negHead) Modify(line []byte) ([]byte, error) {
	if h.n == 0 {
		return nil, io.EOF
	}
	if line == nil {
		return nil, io.EOF
	}

	// Still got room in the buffer, append the current line.
	if len(h.lines) < cap(h.lines) {
		h.lines = append(h.lines, line)
		return nil, nil
	}
	// Insert the new line and pop the first line and return it.
	ret := h.lines[0]
	for i := 0; i < len(h.lines)-1; i++ {
		h.lines[i] = h.lines[i+1]
	}
	h.lines[len(h.lines)-1] = line

	return append(ret, byte('\n')), nil
}

func (h *negHead) Name() string {
	return fmt.Sprintf("head(-%d)", h.n)
}

type tail struct {
	n     int
	lines [][]byte
}

func (t *tail) Modify(line []byte) ([]byte, error) {
	if t.n == 0 {
		return nil, io.EOF
	}
	if line == nil {
		return append(bytes.Join(t.lines, []byte{'\n'}), '\n'), io.EOF
	}

	// Shift all lines and append the new line.
	if len(t.lines) < cap(t.lines) {
		t.lines = append(t.lines, line)
	} else {
		for i := 0; i < len(t.lines)-1; i++ {
			t.lines[i] = t.lines[i+1]
		}
		t.lines[len(t.lines)-1] = line
	}

	return nil, nil
}

func (t *tail) Name() string {
	return fmt.Sprintf("tail(%d)", t.n)
}

type negTail struct {
	n int
}

func (t *negTail) Modify(line []byte) ([]byte, error) {
	if line == nil {
		return nil, io.EOF
	}
	if t.n > 0 {
		t.n--
		return nil, nil
	}
	return append(line, '\n'), nil
}

func (t *negTail) Name() string {
	return fmt.Sprintf("tail(-%d)", t.n)
}
