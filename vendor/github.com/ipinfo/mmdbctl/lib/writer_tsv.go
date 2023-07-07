package lib

import (
	"fmt"
	"io"
	"strings"
)

type TsvWriter struct {
	w io.Writer
}

func NewTsvWriter(w io.Writer) *TsvWriter {
	return &TsvWriter{
		w: w,
	}
}

func (w *TsvWriter) Write(record []string) error {
	_, err := fmt.Fprintln(w.w, strings.Join(record, "\t"))
	return err
}

func (w *TsvWriter) Flush() {
}

func (w *TsvWriter) Error() error {
	return nil
}
