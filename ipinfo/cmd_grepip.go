package main

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
	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

func printHelpGrepIP() {
	fmt.Printf(
		`Usage: %s grepip [<opts>]

Options:
  General:
    --only-matching, -o
      print only matched IP in result line, excluding surrounding content.
    --no-filename, -h
      don't print source of match in result lines when more than 1 source.
    --no-recurse
      don't recurse into more directories in directory sources.
    --help
      show help.

  Filters:
    --ipv4, -4
      match only IPv4 addresses.
    --ipv6, -6
      match only IPv6 addresses.
    --localhost, -l
      match IPs in 127.0.0.0/8 range.
    --bogon, -b
      match any bogon IP.
      full list can be found at https://ipinfo.io/bogon.
`, progBase)
}

func cmdGrepIP() error {
	var fOnlyMatching bool
	var fNoFilename bool
	var fNoRecurse bool
	var fHelp bool
	var fV4 bool
	var fV6 bool
	var fLocalhost bool
	var fBogon bool

	pflag.BoolVarP(&fOnlyMatching, "only-matching", "o", false, "print only matched IPs in result line.")
	pflag.BoolVarP(&fNoFilename, "no-filename", "h", false, "don't print source of match in result lines when more than 1 source.")
	pflag.BoolVarP(&fNoRecurse, "no-recurse", "", false, "don't recurse into more dirs in dir sources.")
	pflag.BoolVarP(&fHelp, "help", "", false, "show help.")
	pflag.BoolVarP(&fV4, "ipv4", "4", false, "print only IPv4 matches.")
	pflag.BoolVarP(&fV6, "ipv6", "6", false, "print only IPv6 matches.")
	pflag.BoolVarP(&fLocalhost, "localhost", "l", false, "match IPs in 127.0.0.0/8 range.")
	pflag.BoolVarP(&fBogon, "bogon", "b", false, "match any bogon IP.")
	pflag.Parse()

	if fHelp {
		printHelpGrepIP()
		return nil
	}

	args := pflag.Args()[1:]

	// require args.
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(args) == 0 && !isStdin {
		printHelpGrepIP()
		return nil
	}

	// if user hasn't forced no-filename, and we have more than 1 source, then
	// output file
	if !fNoFilename && !(len(args) == 0 || (len(args) == 1 && !isStdin)) {
		fNoFilename = false
	} else {
		fNoFilename = true
	}

	// figure out exactly what IP versions we'll match; 0=all, 4=ipv4, 6=ipv6.
	ipv := 0
	if fV4 && fV6 {
		ipv = 0
	} else if fV4 {
		ipv = 4
	} else if fV6 {
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
		start lib.IP6u128
		end   lib.IP6u128
	}
	var bogonRanges4 []iprange4
	var localhostRange4 iprange4
	var bogonRanges6 []iprange6
	var localhostRange6 iprange6
	if fBogon {
		// v4
		bogonRanges4Str := []string{
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
		bogonRanges4 = make([]iprange4, len(bogonRanges4Str))
		for i, bogonRangeStr := range bogonRanges4Str {
			start, end, err := lib.IPRangeStartEndFromCIDR(bogonRangeStr)
			if err != nil {
				panic(err)
			}

			bogonRanges4[i] = iprange4{
				start: start,
				end:   end,
			}
		}

		// v6
		bogonRanges6Str := []string{
			"::/128",
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
		bogonRanges6 = make([]iprange6, len(bogonRanges6Str))
		for i, bogonRangeStr := range bogonRanges6Str {
			start, end, err := lib.IP6RangeStartEndFromCIDR(bogonRangeStr)
			if err != nil {
				panic(err)
			}

			bogonRanges6[i] = iprange6{
				start: start,
				end:   end,
			}
		}
	} else if fLocalhost { // localhost is a subset of bogon
		// v4
		start, end, err := lib.IPRangeStartEndFromCIDR("127.0.0.0/8")
		if err != nil {
			panic(err)
		}

		localhostRange4 = iprange4{
			start: start,
			end:   end,
		}

		// v6
		start6, end6, err := lib.IP6RangeStartEndFromCIDR("::1/128")
		if err != nil {
			panic(err)
		}

		localhostRange6 = iprange6{
			start: start6,
			end:   end6,
		}
	}

	fmtSrc := color.New(color.FgMagenta)
	fmtMatch := color.New(color.Bold, color.FgRed)

	// actual scanner.
	scanrdr := func(src string, r io.Reader) {
		buf := bufio.NewReader(r)

		for {
			d, err := buf.ReadString('\n')
			if err != nil {
				// TODO print error but have a `-q` flag to be quiet.
				return
			}

			// get all matches and then filter.
			var matches [][]int
			allMatches := rexp.FindAllStringIndex(d, -1)
			if !fLocalhost && !fBogon {
				matches = allMatches
			} else {
				matches = make([][]int, 0, len(allMatches))
				for _, m := range allMatches {
					mIPStr := d[m[0]:m[1]]
					mIP := net.ParseIP(mIPStr)
					if strings.Contains(mIPStr, ":") {
						ip := mIP.To16()
						ip128 := lib.IP6u128{
							Hi: binary.BigEndian.Uint64(ip[:8]),
							Lo: binary.BigEndian.Uint64(ip[8:]),
						}

						if fBogon {
							for _, r := range bogonRanges6 {
								if ip128.Gte(r.start) && ip128.Lte(r.end) {
									matches = append(matches, m)
									break
								}
							}
						} else if fLocalhost {
							r := localhostRange6
							if ip128.Gte(r.start) && ip128.Lte(r.end) {
								matches = append(matches, m)
							}
						}
					} else {
						ip := binary.BigEndian.Uint32(mIP.To4())

						if fBogon {
							for _, r := range bogonRanges4 {
								if ip >= r.start && ip <= r.end {
									matches = append(matches, m)
									break
								}
							}
						} else if fLocalhost { // localhost is a subset of bogon
							r := localhostRange4
							if ip >= r.start && ip <= r.end {
								matches = append(matches, m)
							}
						}
					}
				}
			}

			if len(matches) == 0 {
				continue
			}

			// print line up to last match.
			prevMatchEnd := 0
			for _, m := range matches {
				// print source.
				if !fNoFilename && (prevMatchEnd == 0 || fOnlyMatching) {
					fmtSrc.Printf("%s:", src)
				}

				// print pre-match.
				if !fOnlyMatching {
					fmt.Printf("%s", d[prevMatchEnd:m[0]])
				}

				// print match.
				fmtMatch.Printf("%s", d[m[0]:m[1]])
				if fOnlyMatching && prevMatchEnd == 0 && len(matches) > 1 {
					fmt.Printf("\n")
				}

				prevMatchEnd = m[1]
			}

			// print remaining portion and a newline.
			if !fOnlyMatching {
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
				if fNoRecurse && d.IsDir() {
					return fs.SkipDir
				}

				scanfile(path)
				return nil
			})
		}
	}

	return nil
}
