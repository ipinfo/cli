package script

import (
	"fmt"
	"regexp"
)

// Grep filters only line that match the given regexp.
//
// Shell command: `grep <re>`.
func (s Stream) Grep(re *regexp.Regexp) Stream {
	return s.Modify(Grep{Re: re})
}

// Grep is a modifier that filters only line that match `Re`. If Invert was set only line that did
// not match the regex will be returned.
//
// Usage:
//
//  (<Stream object>).Modify(script.Grep{Re: <re>})
//
// Shell command: `grep [-v <Invert>] <Re>`.
type Grep struct {
	Re      *regexp.Regexp
	Inverse bool
}

func (g Grep) Modify(line []byte) (modifed []byte, err error) {
	if line == nil {
		return nil, nil
	}
	if g.Re.Match(line) != g.Inverse {
		return append(line, '\n'), nil
	}
	return nil, nil
}

func (g Grep) Name() string {
	return fmt.Sprintf("grep(%v, invert=%v)", g.Re, g.Inverse)
}
