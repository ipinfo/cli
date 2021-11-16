package lib

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

// CmdGrepIPFlags are flags expected by CmdGrepIP.
type CmdGrepIPFlags struct {
	OnlyMatching bool
	NoFilename   bool
	NoRecurse    bool
	Help         bool
	NoColor      bool
	V4           bool
	V6           bool
	ExclRes      bool
}

// Init initializes the common flags available to CmdGrepIP with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdGrepIPFlags) Init() {
	pflag.BoolVarP(
		&f.OnlyMatching,
		"only-matching", "o", false,
		"print only matched IPs in result line.",
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
		&f.V4,
		"ipv4", "4", false,
		"print only IPv4 matches.",
	)
	pflag.BoolVarP(
		&f.V6,
		"ipv6", "6", false,
		"print only IPv6 matches.",
	)
	pflag.BoolVarP(
		&f.ExclRes,
		"exclude-reserved", "x", false,
		"exclude reserved/bogon IPs.",
	)
}

// CmdGrepIP is the common core logic for the grepip command.
func CmdGrepIP(
	f CmdGrepIPFlags,
	args []string,
	printHelp func(),
) error {
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

	// figure out exactly what IP versions we'll match; 0=all, 4=ipv4, 6=ipv6.
	ipv := 0
	if f.V4 && f.V6 {
		ipv = 0
	} else if f.V4 {
		ipv = 4
	} else if f.V6 {
		ipv = 6
	}

	// prepare regexp
	var rexp *regexp.Regexp
	rexp4 := "((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
	rexp6 := "(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))"
	if ipv == 4 {
		rexp = regexp.MustCompile(rexp4)
	} else if ipv == 6 {
		rexp = regexp.MustCompile(rexp6)
	} else {
		rexp = regexp.MustCompile(rexp4 + "|" + rexp6)
	}

	fmtSrc := color.New(color.FgMagenta)
	fmtMatch := color.New(color.Bold, color.FgRed)

	// actual scanner.
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

				// do one more loop on remaining content.
				hitEOF = true
			} else if err != nil {
				// TODO print error but have a `-q` flag to be quiet.
				return
			}

			// get all matches and then filter.
			var matches [][]int
			allMatches := rexp.FindAllStringIndex(d, -1)
			if !f.ExclRes {
				matches = allMatches
			} else {
				matches = make([][]int, 0, len(allMatches))
				for _, m := range allMatches {
					mIPStr := d[m[0]:m[1]]
					mIP := net.ParseIP(mIPStr)
					if strings.Contains(mIPStr, ":") {
						ip, _ := IP6FromStdIP(mIP.To16())
						for _, r := range bogonIP6List {
							if ip.Gte(r.Start) && ip.Lte(r.End) {
								goto next_match
							}
						}
					} else {
						ip := IPFromStdIP(mIP)
						for _, r := range bogonIP4List {
							if ip >= r.Start && ip <= r.End {
								goto next_match
							}
						}
					}

					matches = append(matches, m)
				next_match:
				}
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
					fmt.Printf("\n")
				}

				prevMatchEnd = m[1]
			}

			// print remaining portion and a newline, if any, and only if we
			// need to print it.
			if !f.OnlyMatching {
				m := matches[len(matches)-1]
				if m[1] < len(d) {
					fmt.Printf("%s", d[m[1]:])
				}
			}
		}
	}

	// opens a file and delegates to scanrdr.
	scanfile := func(path string) {
		f, err := os.Open(path)
		if err != nil {
			// TODO print error but have a `-q` flag to be quiet.
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
