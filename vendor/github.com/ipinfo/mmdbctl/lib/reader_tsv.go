package lib

import (
	"bufio"
	"io"
	"strings"
)

type TsvReader struct {
	r *bufio.Reader
}

func NewTsvReader(r io.Reader) *TsvReader {
	return &TsvReader{
		r: bufio.NewReader(r),
	}
}

func (r *TsvReader) Read() (record []string, err error) {
	line, err := r.r.ReadString('\n')

	// equivalent to perl chomp: remove all extra ending newline content.
	line = strings.TrimRight(line, "\r\n")

	// if this was the end but we have stuff, don't EOF just yet - let the user
	// handle this line, then the underlying reader will EOF on the next read.
	if len(line) > 0 && err == io.EOF {
		err = nil
	}

	return strings.Split(line, "\t"), err
}
