package script

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
)

// To writes the output of the stream to an io.Writer and closes it.
func (s Stream) To(w io.Writer) error {
	var errors *multierror.Error
	if _, err := io.Copy(w, s); err != nil {
		errors = multierror.Append(errors, err)
	}
	if err := s.Close(); err != nil {
		errors = multierror.Append(errors, err)
	}
	return errors.ErrorOrNil()
}

func (s Stream) Iterate(iterator func(line []byte) error) error {
	return s.Modify(ModifyFn(func(line []byte) (modifed []byte, err error) {
		err = iterator(line)
		return nil, err
	})).To(ioutil.Discard)
}

type iterator struct{}

// ToStdout pipes the stdout of the stream to screen.
func (s Stream) ToStdout() error {
	return s.To(os.Stdout)
}

// ToString reads stdout of the stream and returns it as a string.
func (s Stream) ToString() (string, error) {
	var out bytes.Buffer
	err := s.To(&out)
	return out.String(), err

}

// ToFile dumps the output of the stream to a file.
func (s Stream) ToFile(path string) error {
	f, err := File(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return s.To(f)
}

// AppendFile appends the output of the stream to a file.
func (s Stream) AppendFile(path string) error {
	f, err := AppendFile(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return s.To(f)
}

// ToTempFile dumps the output of the stream to a temporary file and returns the temporary files'
// path.
func (s Stream) ToTempFile() (path string, err error) {
	f, err := ioutil.TempFile("", "script-")
	if err != nil {
		return "", err
	}
	defer f.Close()
	return f.Name(), s.To(f)
}

// Discard executes the stream pipeline but discards the output.
func (s Stream) Discard() error {
	return s.To(ioutil.Discard)
}

func File(path string) (io.WriteCloser, error) {
	err := makeDir(path)
	if err != nil {
		return nil, err
	}
	return os.Create(path)
}

func AppendFile(path string) (io.WriteCloser, error) {
	err := makeDir(path)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path); err != nil {
		return File(path)
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0666)
}

func makeDir(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0775)
}
