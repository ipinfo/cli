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
	OnlyMatching  bool
	IncludeCIDRs  bool
	IncludeRanges bool
	CIDRsOnly     bool
	RangesOnly    bool
	NoFilename    bool
	NoRecurse     bool
	Help          bool
	NoColor       bool
	V4            bool
	V6            bool
	ExclRes       bool
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
		&f.IncludeCIDRs,
		"include-cidrs", "", false,
		"print cidrs too.",
	)
	pflag.BoolVarP(
		&f.IncludeRanges,
		"include-ranges", "", false,
		"print ranges too.",
	)
	pflag.BoolVarP(
		&f.CIDRsOnly,
		"cidrs-only", "", false,
		"print cidrs.",
	)
	pflag.BoolVarP(
		&f.RangesOnly,
		"ranges-only", "", false,
		"print ranges.",
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
	if ipv == 4 {
		rexp = ipV4Regex
		if f.CIDRsOnly && f.RangesOnly {
			rexp = v4SubnetRegex
		} else if f.IncludeCIDRs && f.IncludeRanges {
			rexp = v4IpSubnetRegex
		} else if f.IncludeCIDRs {
			rexp = v4IpCidrRegex
		} else if f.IncludeRanges {
			rexp = v4IpRangeRegex
		} else if f.CIDRsOnly {
			rexp = v4CidrRegex
		} else if f.RangesOnly {
			rexp = v4RangeRegex
		}
	} else if ipv == 6 {
		rexp = ipV6Regex
		if f.CIDRsOnly && f.RangesOnly {
			rexp = v6SubnetRegex
		} else if f.IncludeCIDRs && f.IncludeRanges {
			rexp = v6IpSubnetRegex
		} else if f.IncludeCIDRs {
			rexp = v6IpCidrRegex
		} else if f.IncludeRanges {
			rexp = v6IpRangeRegex
		} else if f.CIDRsOnly {
			rexp = v6CidrRegex
		} else if f.RangesOnly {
			rexp = v6RangeRegex
		}
	} else {
		rexp = ipRegex
		if f.CIDRsOnly && f.RangesOnly {
			rexp = subnetRegex
		} else if f.IncludeCIDRs && f.IncludeRanges {
			rexp = ipSubnetRegex
		} else if f.IncludeCIDRs {
			rexp = ipCidrRegex
		} else if f.IncludeRanges {
			rexp = ipRangeRegex
		} else if f.CIDRsOnly {
			rexp = cidrRegex
		} else if f.RangesOnly {
			rexp = rangeRegex
		}
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
					if mIP == nil {
						goto next_match
					}
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
