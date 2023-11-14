package lib

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"golang.org/x/net/idna"
)

// regular expression that matches simple domains and IDN ones.
var domainRegex = regexp.MustCompile(`\b(?:[a-zA-Z0-9-]+\.){1,}[a-zA-Z]{2,}\b|([^\s.]+\.)+[^\s]{2,}\b`)

type CmdGrepDomainFlags struct {
	OnlyMatching bool
	NoFilename   bool
	NoRecurse    bool
	Help         bool
	NoColor      bool
	ExcludePuny  bool
}

func (f *CmdGrepDomainFlags) Init() {
	pflag.BoolVarP(
		&f.OnlyMatching,
		"only-matching", "o", false,
		"print only matched domains in result line.",
	)
	pflag.BoolVarP(
		&f.NoFilename,
		"no-filename", "h", false,
		"don't print source of match in result lines when more than 1 source.",
	)
	pflag.BoolVarP(
		&f.NoRecurse,
		"no-recurse", "", false,
		"don't recurse into more dirs in dir sources.",
	)
	pflag.BoolVarP(
		&f.Help,
		"help", "", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.NoColor,
		"nocolor", "", false,
		"disable color output.",
	)
	pflag.BoolVarP(
		&f.ExcludePuny,
		"no-punycode", "n", false,
		"do not convert domains to punycode.",
	)
}

// CmdGrepDomain is the common core logic for the grepdom command.
func CmdGrepDomain(f CmdGrepDomainFlags, args []string, printHelp func()) error {
	if f.NoColor {
		color.NoColor = true
	}

	if f.Help {
		printHelp()
		return nil
	}

	// require args.
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	srcCnt := len(args)
	if isStdin {
		srcCnt += 1
	}
	if srcCnt == 0 {
		printHelp()
		return nil
	}

	// if user hasn't forced no-filename, and we have more than 1 source, then
	// output file
	if !f.NoFilename && srcCnt > 1 {
		f.NoFilename = false
	} else {
		f.NoFilename = true
	}

	fmtSrc := color.New(color.FgMagenta)
	fmtMatch := color.New(color.Bold, color.FgRed)

	rexp := domainRegex

	// Actual scanner logic for domains
	scanrdr := func(src string, r io.Reader) {
		var hitEOF bool
		buf := bufio.NewReader(r)

		for {
			if hitEOF {
				return
			}

			d, err := buf.ReadString('\n')
			if err == io.EOF {
				if len(d) == 0 {
					return
				}

				// do one more loop for remaining content
				hitEOF = true
			} else if err != nil {
				return
			}

			var matches [][]int
			allMatches := rexp.FindAllStringIndex(d, -1)
			if !f.ExcludePuny {
				for _, m := range allMatches {
					mDomainStr := d[m[0]:m[1]]
					puny, err := idna.ToASCII(mDomainStr)
					if err != nil {
						goto next_match
					}
					// print original domain and the corresponding punycode
					fmt.Printf("Original domain: %s\n", mDomainStr)
					fmt.Printf("Punycode: %s\n", puny)
					matches = append(matches, m)
				next_match:
				}
			} else {
				matches = allMatches
			}

			// no match?
			if len(matches) == 0 {
				continue
			}

			// print line up to last match.
			prevMatchEnd := 0
			for _, m := range matches {
				// print source if requested, but only if we're printing
				// 1 match per line, or this is the first match of the line.
				if !f.NoFilename && (prevMatchEnd == 0 || f.OnlyMatching) {
					fmtSrc.Printf("%s:", src)
				}

				// print everything up to the current match.
				if !f.OnlyMatching {
					fmt.Printf("%s", d[prevMatchEnd:m[0]])
				}
				// print the match itself.
				fmtMatch.Printf("%s", d[m[0]:m[1]])

				if f.OnlyMatching {
					fmt.Printf("\n\n")
				}

				prevMatchEnd = m[1]
			}

			// print remaining portion and a newline, if any, and only if we
			// need to print it.
			if !f.OnlyMatching {
				m := matches[len(matches)-1]
				if m[1] < len(d) {
					fmt.Printf("%s\n", d[m[1]:])
				}
			}
		}
	}

	// opens a file and delegates to scanrdr.
	scanfile := func(path string) {
		f, err := os.Open(path)
		if err != nil {
			return
		}

		scanrdr(path, f)
	}

	// scan stdin first.
	if isStdin {
		scanrdr("(standard input)", os.Stdin)
	}

	// scan all args.
	for _, arg := range args {
		fi, err := os.Stat(arg)
		if err != nil {
			continue
		}

		switch mode := fi.Mode(); {
		case mode.IsRegular():
			scanfile(arg)
		case mode.IsDir():
			filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				// skip input dir.
				if arg == path {
					return nil
				}
				// don't recurse if requested.
				if f.NoRecurse && d.IsDir() {
					return fs.SkipDir
				}

				scanfile(path)
				return nil
			})
		}
	}

	return nil
}
