package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
)

var progBase = filepath.Base(os.Args[0])
var version = "2.6.0"

// global flags.
var fHelp bool
var fNoCache bool
var fNoColor bool

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
	case lib.StrIsIPStr(cmd):
		err = cmdIP(cmd)
	case lib.StrIsASNStr(cmd):
		asn := strings.ToUpper(cmd)
		err = cmdASN(asn)
	case len(cmd) >= 3 && strings.IndexByte(cmd, '.') != -1:
		err = cmdDomain(cmd)
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
	case cmd == "cidr2ip":
		err = cmdCIDR2IP()
	case cmd == "range2cidr":
		err = cmdRange2CIDR()
	case cmd == "range2ip":
		err = cmdRange2IP()
	case cmd == "randip":
		err = cmdRandIP()
	case cmd == "cache":
		err = cmdCache()
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
