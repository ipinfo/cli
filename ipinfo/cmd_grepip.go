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
      print only IPv4 matches.
    --ipv6, -6
      print only IPv6 matches.
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

	// each range is a 2-tuple of start/end IPs of the range, inclusive.
	var bogonRanges [][]uint32
	var localhostRange []uint32

	rexp := regexp.MustCompile("[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}")
	fmtSrc := color.New(color.FgMagenta)
	fmtMatch := color.New(color.Bold, color.FgRed)
	if fBogon {
		bogonRangesStr := []string{
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
		bogonRanges = make([][]uint32, len(bogonRangesStr))
		for i, bogonRangeStr := range bogonRangesStr {
			start, end, err := lib.IPRangeStartEndFromCIDR(bogonRangeStr)
			if err != nil {
				panic(err)
			}

			bogonRanges[i] = []uint32{start, end}
		}
	}
	if fLocalhost {
		start, end, err := lib.IPRangeStartEndFromCIDR("127.0.0.0/8")
		if err != nil {
			panic(err)
		}

		localhostRange = []uint32{start, end}
	}

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
					mIP := net.ParseIP(d[m[0]:m[1]])
					ip := binary.BigEndian.Uint32(mIP.To4())

					if fBogon {
						for _, bogonRange := range bogonRanges {
							if ip >= bogonRange[0] && ip <= bogonRange[1] {
								matches = append(matches, m)
								break
							}
						}
					} else if fLocalhost { // localhost is a subset of bogon
						if ip >= localhostRange[0] && ip <= localhostRange[1] {
							matches = append(matches, m)
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

	// if user hasn't forced no-filename, and we have more than 1 source, then
	// output file
	if !fNoFilename && !(len(args) == 0 || (len(args) == 1 && !isStdin)) {
		fNoFilename = false
	} else {
		fNoFilename = true
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
