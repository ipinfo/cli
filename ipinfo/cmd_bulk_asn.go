package main

import (
	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

// TODO:
// Figure out help messages

func cmdBulkASN() error {
	f := lib.CmdBulkASNFlags{}
	f.Init()
	pflag.Parse()

	ii = prepareIpinfoClient(f.Token)

	data, err := lib.CmdBulkASN(ii, pflag.Args())
	if err != nil {
		return err
	}

	if len(f.Field) > 0 {
		return outputFieldBatchASNDetails(data, f.Field, true, true)
	}

	if f.Yaml {
		return outputYAML(data)
	}

	return outputJSON(data)
}
