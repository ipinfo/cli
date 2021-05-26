package script

import (
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-multierror"
)

// Cat outputs the contents of the given files.
//
// Shell command: cat <path>.
func Cat(paths ...string) Stream {
	var (
		readers []io.Reader
		closers multicloser
		errors  *multierror.Error
	)

	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("open path %s: %v", path, err))
		} else {
			readers = append(readers, f)
			closers = append(closers, f)
		}
	}

	return Stream{
		r:     readcloser{Reader: io.MultiReader(readers...), Closer: closers},
		stage: "cat",
		err:   errors.ErrorOrNil(),
	}

}

type multicloser []io.Closer

func (mc multicloser) Close() error {
	var errors *multierror.Error
	for _, c := range mc {
		if err := c.Close(); err != nil {
			errors = multierror.Append(errors, err)
		}
	}
	return errors.ErrorOrNil()
}
