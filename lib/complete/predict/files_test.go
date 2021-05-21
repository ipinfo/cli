package predict

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestFiles(t *testing.T) {
	err := os.Chdir("testdata")
	if err != nil {
		panic(err)
	}
	defer os.Chdir("..")

	tests := []struct {
		name     string
		p        FilesPredictor
		prefixes []string
		want     []string
	}{
		{
			name:     "files/txt",
			p:        Files("*.txt"),
			prefixes: []string{""},
			want:     []string{"./", "dir/", "outer/", "a.txt", "b.txt", "c.txt", ".dot.txt"},
		},
		{
			name:     "files/txt",
			p:        Files("*.txt"),
			prefixes: []string{"./dir/"},
			want:     []string{"./dir/"},
		},
		{
			name:     "complete files inside dir if it is the only match",
			p:        Files("foo"),
			prefixes: []string{"./dir/", "./d"},
			want:     []string{"./dir/", "./dir/foo"},
		},
		{
			name:     "complete files inside dir when argList includes file name",
			p:        Files("*"),
			prefixes: []string{"./dir/f", "./dir/foo"},
			want:     []string{"./dir/foo"},
		},
		{
			name:     "files/md",
			p:        Files("*.md"),
			prefixes: []string{""},
			want:     []string{"./", "dir/", "outer/", "readme.md"},
		},
		{
			name:     "files/md with ./ prefix",
			p:        Files("*.md"),
			prefixes: []string{".", "./"},
			want:     []string{"./", "./dir/", "./outer/", "./readme.md"},
		},
		{
			name:     "dirs",
			p:        Dirs("*"),
			prefixes: []string{"di", "dir", "dir/"},
			want:     []string{"dir/"},
		},
		{
			name:     "dirs with ./ prefix",
			p:        Dirs("*"),
			prefixes: []string{"./di", "./dir", "./dir/"},
			want:     []string{"./dir/"},
		},
		{
			name:     "predict anything in dir",
			p:        Files("*"),
			prefixes: []string{"dir", "dir/", "di"},
			want:     []string{"dir/", "dir/foo", "dir/bar"},
		},
		{
			name:     "predict anything in dir with ./ prefix",
			p:        Files("*"),
			prefixes: []string{"./dir", "./dir/", "./di"},
			want:     []string{"./dir/", "./dir/foo", "./dir/bar"},
		},
		{
			name:     "root directories",
			p:        Dirs("*"),
			prefixes: []string{""},
			want:     []string{"./", "dir/", "outer/"},
		},
		{
			name:     "root directories with ./ prefix",
			p:        Dirs("*"),
			prefixes: []string{".", "./"},
			want:     []string{"./", "./dir/", "./outer/"},
		},
		{
			name:     "nested directories",
			p:        Dirs("*.md"),
			prefixes: []string{"ou", "outer", "outer/"},
			want:     []string{"outer/", "outer/inner/"},
		},
		{
			name:     "nested directories with ./ prefix",
			p:        Dirs("*.md"),
			prefixes: []string{"./ou", "./outer", "./outer/"},
			want:     []string{"./outer/", "./outer/inner/"},
		},
		{
			name:     "nested inner directory",
			p:        Files("*.md"),
			prefixes: []string{"outer/i"},
			want:     []string{"outer/inner/", "outer/inner/readme.md"},
		},
	}

	for _, tt := range tests {
		for _, prefix := range tt.prefixes {
			t.Run(tt.name+"/prefix="+prefix, func(t *testing.T) {

				matches := tt.p.Predict(prefix)

				sort.Strings(matches)
				sort.Strings(tt.want)

				got := strings.Join(matches, ",")
				want := strings.Join(tt.want, ",")

				if got != want {
					t.Errorf("failed %s\ngot = %s\nwant: %s", t.Name(), got, want)
				}
			})
		}
	}
}

func TestMatchFile(t *testing.T) {
	// Change to tests directory for testing completion of
	// files and directories
	err := os.Chdir("testdata")
	if err != nil {
		panic(err)
	}
	defer os.Chdir("..")

	type matcherTest struct {
		prefix string
		want   bool
	}

	tests := []struct {
		long  string
		tests []matcherTest
	}{
		{
			long: "file.txt",
			tests: []matcherTest{
				{prefix: "", want: true},
				{prefix: "f", want: true},
				{prefix: "./f", want: true},
				{prefix: "./.", want: false},
				{prefix: "file.", want: true},
				{prefix: "./file.", want: true},
				{prefix: "file.txt", want: true},
				{prefix: "./file.txt", want: true},
				{prefix: "other.txt", want: false},
				{prefix: "/other.txt", want: false},
				{prefix: "/file.txt", want: false},
				{prefix: "/fil", want: false},
				{prefix: "/file.txt2", want: false},
				{prefix: "/.", want: false},
			},
		},
		{
			long: "./file.txt",
			tests: []matcherTest{
				{prefix: "", want: true},
				{prefix: "f", want: true},
				{prefix: "./f", want: true},
				{prefix: "./.", want: false},
				{prefix: "file.", want: true},
				{prefix: "./file.", want: true},
				{prefix: "file.txt", want: true},
				{prefix: "./file.txt", want: true},
				{prefix: "other.txt", want: false},
				{prefix: "/other.txt", want: false},
				{prefix: "/file.txt", want: false},
				{prefix: "/fil", want: false},
				{prefix: "/file.txt2", want: false},
				{prefix: "/.", want: false},
			},
		},
		{
			long: "/file.txt",
			tests: []matcherTest{
				{prefix: "", want: true},
				{prefix: "f", want: false},
				{prefix: "./f", want: false},
				{prefix: "./.", want: false},
				{prefix: "file.", want: false},
				{prefix: "./file.", want: false},
				{prefix: "file.txt", want: false},
				{prefix: "./file.txt", want: false},
				{prefix: "other.txt", want: false},
				{prefix: "/other.txt", want: false},
				{prefix: "/file.txt", want: true},
				{prefix: "/fil", want: true},
				{prefix: "/file.txt2", want: false},
				{prefix: "/.", want: false},
			},
		},
		{
			long: "./",
			tests: []matcherTest{
				{prefix: "", want: true},
				{prefix: ".", want: true},
				{prefix: "./", want: true},
				{prefix: "./.", want: false},
			},
		},
	}

	for _, tt := range tests {
		for _, ttt := range tt.tests {
			name := fmt.Sprintf("long=%q&prefix=%q", tt.long, ttt.prefix)
			t.Run(name, func(t *testing.T) {
				got := matchFile(tt.long, ttt.prefix)
				if got != ttt.want {
					t.Errorf("Failed %s: got = %t, want: %t", name, got, ttt.want)
				}
			})
		}
	}
}
