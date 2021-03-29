package main

import (
	"fmt"
	"os"
	"path/filepath"
	"io"
	"io/fs"

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
      don't print source of match in result lines.
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
	pflag.BoolVarP(&fNoFilename, "no-filename", "h", false, "don't print source of match in result lines.")
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
	if len(args) == 0 && (stat.Mode() & os.ModeCharDevice) != 0 {
		printHelpGrepIP()
		return nil
	}

	// actual scanner.
	scanrdr := func(src string, r io.Reader) {
		fmt.Printf("scanning %v\n", src)
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
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		scanrdr("(stdin)", os.Stdin)
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
