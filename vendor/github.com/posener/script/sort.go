package script

import (
	"bufio"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
)

// Sort returns a stream with lines ordered alphabetically.
//
// Shell command: `wc`.
func (s Stream) Sort(reverse bool) Stream {
	var (
		lines  []string
		errors *multierror.Error
	)
	scanner := bufio.NewScanner(s.r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		errors = multierror.Append(errors, fmt.Errorf("scanning stream: %s", err))
	}

	sort.Slice(lines, func(i, j int) bool { return (lines[i] < lines[j]) != reverse })

	var out strings.Builder
	for _, line := range lines {
		_, err := out.WriteString(line + "\n")
		if err != nil {
			errors = multierror.Append(errors, fmt.Errorf("writing line %q: %s", line, err))
		}
	}

	return Stream{
		stage: fmt.Sprintf("sort(%v)", reverse),
		r:     strings.NewReader(out.String()),
		err:   errors.ErrorOrNil(),
	}
}
