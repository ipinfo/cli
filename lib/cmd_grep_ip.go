package lib

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

// CompletionsGrepIP are the completions for the grepip command.
var CompletionsGrepIP = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-o":                 predict.Nothing,
		"--only-matching":    predict.Nothing,
		"-h":                 predict.Nothing,
		"--no-filename":      predict.Nothing,
		"--no-recurse":       predict.Nothing,
		"--help":             predict.Nothing,
		"--nocolor":          predict.Nothing,
		"-4":                 predict.Nothing,
		"--ipv4":             predict.Nothing,
		"-6":                 predict.Nothing,
		"--ipv6":             predict.Nothing,
		"-x":                 predict.Nothing,
		"--exclude-reserved": predict.Nothing,
	},
}

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

// CmdGrepIP is the common core logic for the `grepip` command-line utility.
func CmdGrepIP(f CmdGrepIPFlags, args []string, printHelp func()) error {
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
	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}

	// if user hasn't forced no-filename, and we have more than 1 source, then
	// output file
	if !f.NoFilename && !(len(args) == 0 || (len(args) == 1 && !isStdin)) {
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
	rexp4 := "([0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3})"
	rexp6 := "(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))"
	if ipv == 4 {
		rexp = regexp.MustCompile(rexp4)
	} else if ipv == 6 {
		rexp = regexp.MustCompile(rexp6)
	} else {
		rexp = regexp.MustCompile(rexp4 + "|" + rexp6)
	}

	// prepare bogon/localhost ranges
	type iprange4 struct {
		start uint32
		end   uint32
	}
	type iprange6 struct {
		start IP6u128
		end   IP6u128
	}
	var exclRanges4 []iprange4
	var exclRanges6 []iprange6
	if f.ExclRes {
		// v4
		exclRanges4Str := []string{
			"0.0.0.0/8",
			"10.0.0.0/8",
			"100.64.0.0/10",
			"127.0.0.0/8",
			"169.254.0.0/16",
			"172.16.0.0/12",
			"192.0.0.0/24",
			"192.0.2.0/24",
			"192.168.0.0/16",
			"198.18.0.0/15",
			"198.51.100.0/24",
			"203.0.113.0/24",
			"224.0.0.0/4",
			"240.0.0.0/4",
			"255.255.255.255/32",
		}
		exclRanges4 = make([]iprange4, len(exclRanges4Str))
		for i, bogonRangeStr := range exclRanges4Str {
			start, end, err := IPRangeStartEndFromCIDR(bogonRangeStr)
			if err != nil {
				panic(err)
			}

			exclRanges4[i] = iprange4{
				start: start,
				end:   end,
			}
		}

		// v6
		exclRanges6Str := []string{
			"::/128",
			"::1/128",
			"::ffff:0:0/96",
			"::/96",
			"100::/64",
			"2001:10::/28",
			"2001:db8::/32",
			"fc00::/7",
			"fe80::/10",
			"fec0::/10",
			"ff00::/8",
			// 6to4 bogon ranges
			"2002::/24",
			"2002:a00::/24",
			"2002:7f00::/24",
			"2002:a9fe::/32",
			"2002:ac10::/28",
			"2002:c000::/40",
			"2002:c000:200::/40",
			"2002:c0a8::/32",
			"2002:c612::/31",
			"2002:c633:6400::/40",
			"2002:cb00:7100::/40",
			"2002:e000::/20",
			"2002:f000::/20",
			"2002:ffff:ffff::/48",
			// teredo
			"2001::/40",
			"2001:0:a00::/40",
			"2001:0:7f00::/40",
			"2001:0:a9fe::/48",
			"2001:0:ac10::/44",
			"2001:0:c000::/56",
			"2001:0:c000:200::/56",
			"2001:0:c0a8::/48",
			"2001:0:c612::/47",
			"2001:0:c633:6400::/56",
			"2001:0:cb00:7100::/56",
			"2001:0:e000::/36",
			"2001:0:f000::/36",
			"2001:0:ffff:ffff::/64",
		}
		exclRanges6 = make([]iprange6, len(exclRanges6Str))
		for i, bogonRangeStr := range exclRanges6Str {
			start, end, err := IP6RangeStartEndFromCIDR(bogonRangeStr)
			if err != nil {
				panic(err)
			}

			exclRanges6[i] = iprange6{
				start: start,
				end:   end,
			}
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
			if err == io.EOF && len(d) > 0 {
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
						ip := mIP.To16()
						ip128 := IP6u128{
							Hi: binary.BigEndian.Uint64(ip[:8]),
							Lo: binary.BigEndian.Uint64(ip[8:]),
						}

						for _, r := range exclRanges6 {
							if ip128.Gte(r.start) && ip128.Lte(r.end) {
								goto next_match
							}
						}
					} else {
						ip := binary.BigEndian.Uint32(mIP.To4())

						for _, r := range exclRanges4 {
							if ip >= r.start && ip <= r.end {
								goto next_match
							}
						}
					}

					matches = append(matches, m)
				next_match:
				}
			}

			if len(matches) == 0 {
				continue
			}

			// print line up to last match.
			prevMatchEnd := 0
			for _, m := range matches {
				// print source.
				if !f.NoFilename && (prevMatchEnd == 0 || f.OnlyMatching) {
					fmtSrc.Printf("%s:", src)
				}

				// print pre-match.
				if !f.OnlyMatching {
					fmt.Printf("%s", d[prevMatchEnd:m[0]])
				}

				// print match.
				fmtMatch.Printf("%s", d[m[0]:m[1]])
				if f.OnlyMatching && prevMatchEnd == 0 && len(matches) > 1 {
					fmt.Printf("\n")
				}

				prevMatchEnd = m[1]
			}

			// print remaining portion and a newline.
			if !f.OnlyMatching {
				m := matches[len(matches)-1]
				if m[1] < len(d) {
					fmt.Printf("%s", d[m[1]:len(d)-1])
				}
			}
			fmt.Printf("\n")
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
