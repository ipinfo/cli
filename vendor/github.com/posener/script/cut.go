package script

import (
	"bytes"
	"fmt"
)

// Cut takes selected fields from each line. The fields are 1 based (first field is 1).
//
// Shell command: `cut -f<Fields>`.
func (s Stream) Cut(fields ...int) Stream {
	return s.Modify(Cut{Fields: fields})
}

// Cut is a `Modifier` that takes selected fields from each line according to a given delimiter.
// The default delimiter is tab.
//
// Shell command: `cut -d<Delim> -f<Fields>`.
type Cut struct {
	// Fields defines which fields will be collected to the output of the command. The fields are 1
	// based (first field is 1).
	Fields []int
	// Delim is the delimited by which the fields of each line are sparated.
	Delim []byte
}

func (c Cut) Modify(line []byte) (modifed []byte, err error) {
	if line == nil {
		return nil, nil
	}
	if len(c.Fields) == 0 {
		return nil, nil
	}
	if len(c.Delim) == 0 {
		c.Delim = []byte{'\t'}
	}

	parts := bytes.Split(line, c.Delim)
	out := make([][]byte, 0, len(parts))
	for _, i := range c.Fields {
		i-- // Fields are 1 based, translate to zero base.
		if i < len(parts) {
			out = append(out, parts[i])
		}
	}
	return append(bytes.Join(out, c.Delim), '\n'), nil
}

func (c Cut) Name() string {
	return fmt.Sprintf("cut(%v, delim=%v)", c.Fields, c.Delim)
}
