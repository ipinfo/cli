package script

import (
	"bufio"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Sort returns a stream with lines ordered alphabetically.
//
// Shell command: `wc`.
func (s Stream) Sort(reverse bool) Stream {
	var (
		lines []string
		merr  error
	)
	scanner := bufio.NewScanner(s.r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		merr = errors.Join(merr, fmt.Errorf("scanning stream: %w", err))
	}

	sort.Slice(lines, func(i, j int) bool { return (lines[i] < lines[j]) != reverse })

	var out strings.Builder
	for _, line := range lines {
		_, err := out.WriteString(line + "\n")
		if err != nil {
			merr = errors.Join(merr, fmt.Errorf("writing line %q: %w", line, err))
		}
	}

	return Stream{
		stage: fmt.Sprintf("sort(%v)", reverse),
		r:     strings.NewReader(out.String()),
		err:   merr,
	}
}
