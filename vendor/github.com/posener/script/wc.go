package script

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

// Count represents the output of `wc` shell command.
type Count struct {
	// Stream can be used to pipe the output of wc.
	Stream
	// Count the number of lines, words and chars in the input.
	Lines, Words, Chars int
}

// Wc counts the number of lines, words and characters.
//
// Shell command: `wc`.
func (s Stream) Wc() Count {
	defer s.Close()

	var (
		count Count
		merr  error
	)
	scanner := bufio.NewScanner(s)
	for scanner.Scan() {
		count.Lines++
		count.Chars += len(scanner.Text()) + 1
		count.Words += countWords(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		merr = errors.Join(merr, fmt.Errorf("scanning stream: %w", err))
	}

	count.Stream = Stream{
		stage: "wc",
		r:     strings.NewReader(count.String()),
		err:   merr,
	}
	return count
}

func (c Count) String() string {
	return fmt.Sprintf("%d\t%d\t%d\n", c.Lines, c.Words, c.Chars)
}

func countWords(s string) int {
	// TODO: improve performance.
	return len(strings.Fields(s))
}
