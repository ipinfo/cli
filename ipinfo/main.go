package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/ipinfo/complete/v3"
	"github.com/ipinfo/complete/v3/predict"
)

var progBase = filepath.Base(os.Args[0])
var version = "1.1.5"

var ii *ipinfo.Client

func prepareIpinfoClient(tok string) *ipinfo.Client {
	var _ii *ipinfo.Client

	if tok == "" {
		tok, _ = restoreToken()
	}

	_ii = ipinfo.NewClient(nil, nil, tok)
	_ii.UserAgent = fmt.Sprintf(
		"IPinfoCli/%s (os/%s - arch/%s)",
		version, runtime.GOOS, runtime.GOARCH,
	)
	return _ii
}

func main() {
	var err error
	var cmd string

	// TODO split each subcommand's subtree to its own file and use it here.
	completions := &complete.Command{
		Sub: map[string]*complete.Command{
			"myip": &complete.Command{
				Flags: map[string]complete.Predictor{
					"-h": predict.Nothing,
					"--help": predict.Nothing,
				},
			},
			"bulk": &complete.Command{
				Flags: map[string]complete.Predictor{
					"-t": predict.Nothing,
					"--token": predict.Nothing,
					"-h": predict.Nothing,
					"--help": predict.Nothing,
					"-f": predict.Nothing,
					"--field": predict.Nothing,
					"--nocolor": predict.Nothing,
					"-j": predict.Nothing,
					"--json": predict.Nothing,
					"-c": predict.Nothing,
					"--csv": predict.Nothing,
				},
			},
			"summarize": &complete.Command{
				Flags: map[string]complete.Predictor{
					"-t": predict.Nothing,
					"--token": predict.Nothing,
					"-h": predict.Nothing,
					"--help": predict.Nothing,
					"--nocolor": predict.Nothing,
					"-p": predict.Nothing,
					"--pretty": predict.Nothing,
					"-j": predict.Nothing,
					"--json": predict.Nothing,
				},
			},
			"map": &complete.Command{
				Flags: map[string]complete.Predictor{
					"-h": predict.Nothing,
					"--help": predict.Nothing,
				},
			},
			"prips": &complete.Command{
				Flags: map[string]complete.Predictor{
					"-h": predict.Nothing,
					"--help": predict.Nothing,
				},
			},
			"grepip": &complete.Command{
				Flags: map[string]complete.Predictor{
					"-o": predict.Nothing,
					"--only-matching": predict.Nothing,
					"-h": predict.Nothing,
					"--no-filename": predict.Nothing,
					"--no-recurse": predict.Nothing,
					"--help": predict.Nothing,
					"--nocolor": predict.Nothing,
					"-4": predict.Nothing,
					"--ipv4": predict.Nothing,
					"-6": predict.Nothing,
					"--ipv6": predict.Nothing,
					"-x": predict.Nothing,
					"--exclude-reserved": predict.Nothing,
				},
			},
			"login": &complete.Command{
				Flags: map[string]complete.Predictor{
					"-t": predict.Nothing,
					"--token": predict.Nothing,
					"-h": predict.Nothing,
					"--help": predict.Nothing,
				},
			},
			"logout": &complete.Command{
				Flags: map[string]complete.Predictor{
					"-h": predict.Nothing,
					"--help": predict.Nothing,
				},
			},
			"version": &complete.Command{},
		},
	}
	completions.Complete(progBase)

	// obey NO_COLOR env var.
	if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch {
	case lib.IsIP(cmd):
		err = cmdIP(cmd)
	case lib.IsASN(cmd):
		asn := strings.ToUpper(cmd)
		err = cmdASN(asn)
	case cmd == "myip":
		err = cmdMyIP()
	case cmd == "bulk":
		err = cmdBulk()
	case cmd == "summarize" || cmd == "sum":
		err = cmdSum()
	case cmd == "map":
		err = cmdMap()
	case cmd == "prips":
		err = cmdPrips()
	case cmd == "grepip":
		err = cmdGrepIP()
	case cmd == "login":
		err = cmdLogin()
	case cmd == "logout":
		err = cmdLogout()
	case cmd == "version" || cmd == "vsn" || cmd == "v":
		err = cmdVersion()
	default:
		err = cmdDefault()
	}

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
