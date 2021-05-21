package complete

import (
	"sort"
	"testing"
)

// Test is a testing helper function for testing bash completion of a given
// completer.
func Test(t *testing.T, cmp Completer, args string, want []string) {
	t.Helper()
	got, err := completer{Completer: cmp, args: Parse(args)}.complete()
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(got)
	sort.Strings(want)
	if len(want) != len(got) {
		t.Errorf("got != want: want = %+v, got = %+v", want, got)
		return
	}
	for i := range want {
		if want[i] != got[i] {
			t.Errorf("got != want: want = %+v, got = %+v", want, got)
			return
		}
	}
}
