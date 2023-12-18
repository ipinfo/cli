package script

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Cat outputs the contents of the given files.
//
// Shell command: cat <path>.
func Cat(paths ...string) Stream {
	var (
		readers []io.Reader
		closers multicloser
		merr    error
	)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			merr = errors.Join(merr, fmt.Errorf("open path %s: %w", path, err))
		} else {
			readers = append(readers, f)
			closers = append(closers, f)
		}
	}

	return Stream{
		r:     readcloser{Reader: io.MultiReader(readers...), Closer: closers},
		stage: "cat",
		err:   merr,
	}
}

type multicloser []io.Closer

func (mc multicloser) Close() error {
	var merr error
	for _, c := range mc {
		if err := c.Close(); err != nil {
			merr = errors.Join(merr, err)
		}
	}
	return merr
}
