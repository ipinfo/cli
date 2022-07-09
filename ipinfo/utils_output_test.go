package main

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ipinfo/go/v2/ipinfo"
)

func TestCountryOutput(t *testing.T) {
	re := regexp.MustCompile(`(?U)- Country\W+([A-Za-z ]+(\s\([A-Z]{2}\))?)\n`)
	tests := []struct {
		name string
		have *ipinfo.Core
		want string
	}{
		{
			name: "valid without missing country",
			have: &ipinfo.Core{
				CountryName: "United States",
				Country:     "US",
			},
			want: "United States (US)",
		},
		{
			name: "missing country no empty ()",
			have: &ipinfo.Core{
				CountryName: "United States",
				Country:     "",
			},
			want: "United States",
		},
		{
			name: "totally missing country",
			have: &ipinfo.Core{
				CountryName: "",
				Country:     "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			outputFriendlyCore(tt.have)
			w.Close()

			var buf bytes.Buffer
			io.Copy(&buf, r)
			out := buf.Bytes()

			match := re.FindSubmatch(out)
			have := strings.TrimLeft(string(match[1]), " ")

			if have != tt.want {
				t.Errorf("outputFriendlyCore() = %q, want %q", have, tt.want)
			}

			os.Stdout = old
		})
	}
}
