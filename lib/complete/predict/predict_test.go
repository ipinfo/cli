package predict

import (
	"testing"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/stretchr/testify/assert"
)

func TestPredict(t *testing.T) {
	tests := []struct {
		name   string
		p      complete.Predictor
		prefix string
		want   []string
	}{
		{
			name: "set",
			p:    Set{"a", "b", "c"},
			want: []string{"a", "b", "c"},
		},
		{
			name: "set/empty",
			p:    Set{},
			want: []string{},
		},
		{
			name: "or: word with nil",
			p:    Or(Set{"a"}, nil),
			want: []string{"a"},
		},
		{
			name: "or: nil with word",
			p:    Or(nil, Set{"a"}),
			want: []string{"a"},
		},
		{
			name: "or: word with word with word",
			p:    Or(Set{"a"}, Set{"b"}, Set{"c"}),
			want: []string{"a", "b", "c"},
		},
		{
			name: "something",
			p:    Something,
			want: []string{""},
		},
		{
			name:   "nothing",
			p:      Nothing,
			prefix: "a",
			want:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.Predict(tt.prefix)
			assert.ElementsMatch(t, tt.want, got, "Got: %+v", got)
		})
	}
}
