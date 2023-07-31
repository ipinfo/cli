package main

import (
	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
	"os"
)

func cmdBulkIpAsn() error {
	f := lib.CmdBulkIpAsnFlags{}
	f.Init()
	pflag.Parse()
	return lib.CmdBulkIpAsn(f, os.Args[1:])
}
