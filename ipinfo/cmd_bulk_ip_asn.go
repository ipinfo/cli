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

	ii = prepareIpinfoClient(f.Token)

	return lib.CmdBulkIpAsn(f, ii, os.Args[1:])
}
