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
)

var progBase = filepath.Base(os.Args[0])
var version = "2.0.0"

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

	// obey NO_COLOR env var.
	if os.Getenv("NO_COLOR") != "" {
		color.NoColor = true
	}

	handleCompletions()

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch {
	case lib.StrIsIP(cmd):
		err = cmdIP(cmd)
	case lib.StrIsASN(cmd):
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
	case cmd == "cidr2range":
		err = cmdCIDR2Range()
	case cmd == "range2cidr":
		err = cmdRange2CIDR()
	case cmd == "login":
		err = cmdLogin()
	case cmd == "logout":
		err = cmdLogout()
	case cmd == "completion":
		err = cmdCompletion()
	case cmd == "version" || cmd == "vsn" || cmd == "v":
		err = cmdVersion()
	default:
		err = cmdDefault()
	}

	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
}
