package lib

import (
	"reflect"
	"testing"
)

func TestGrepIPMatches(t *testing.T) {
	cases := []struct {
		name  string
		line  string
		flags CmdGrepIPFlags
		want  []string
	}{
		{
			name:  "public ip kept",
			line:  "115.176.70.246",
			flags: CmdGrepIPFlags{V4: true, IncludeCIDRs: true, ExclRes: true},
			want:  []string{"115.176.70.246"},
		},
		{
			name:  "public cidr kept",
			line:  "103.221.118.124/14",
			flags: CmdGrepIPFlags{V4: true, IncludeCIDRs: true, ExclRes: true},
			want:  []string{"103.221.118.124/14"},
		},
		{
			name:  "bogon cidr filtered",
			line:  "10.0.0.0/24",
			flags: CmdGrepIPFlags{V4: true, IncludeCIDRs: true, ExclRes: true},
			want:  []string{},
		},
		{
			name:  "bogon range filtered",
			line:  "10.0.0.0-10.0.0.5",
			flags: CmdGrepIPFlags{V4: true, IncludeRanges: true, ExclRes: true},
			want:  []string{},
		},
		{
			name:  "public range kept",
			line:  "8.8.8.0-8.8.8.255",
			flags: CmdGrepIPFlags{V4: true, IncludeRanges: true, ExclRes: true},
			want:  []string{"8.8.8.0-8.8.8.255"},
		},

		{
			name:  "public v6 cidr kept",
			line:  "2606:4700::/32",
			flags: CmdGrepIPFlags{V6: true, IncludeCIDRs: true, ExclRes: true},
			want:  []string{"2606:4700::/32"},
		},
		{
			name:  "link-local v6 cidr filtered",
			line:  "fe80::/10",
			flags: CmdGrepIPFlags{V6: true, IncludeCIDRs: true, ExclRes: true},
			want:  []string{},
		},
		{
			name:  "public v6 range kept",
			line:  "2606:4700::1-2606:4700::5",
			flags: CmdGrepIPFlags{V6: true, IncludeRanges: true, ExclRes: true},
			want:  []string{"2606:4700::1-2606:4700::5"},
		},
		{
			name:  "link-local v6 range filtered",
			line:  "fe80::1-fe80::5",
			flags: CmdGrepIPFlags{V6: true, IncludeRanges: true, ExclRes: true},
			want:  []string{},
		},

		{
			name:  "no exclude-reserved keeps everything",
			line:  "10.0.0.0/24",
			flags: CmdGrepIPFlags{V4: true, IncludeCIDRs: true},
			want:  []string{"10.0.0.0/24"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := GrepIPMatches(tc.line, tc.flags)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("GrepIPMatches(%q) = %v, want %v", tc.line, got, tc.want)
			}
		})
	}
}
