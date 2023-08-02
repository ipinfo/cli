package main

import (
	"github.com/ipinfo/cli/lib"
	"github.com/spf13/pflag"
)

func cmdBulkIP() error {
	f := lib.CmdBulkIPFlags{}
	f.Init()
	pflag.Parse()

	ii = prepareIpinfoClient(f.Token)

	data, err := lib.CmdBulkIP(ii, pflag.Args())
	if err != nil {
		return err
	}

	if len(f.Field) > 0 {
		return outputFieldBatchCore(data, f.Field, true, true)
	}

	if f.Csv {
		return outputCSVBatchCore(data)
	}
	if f.Yaml {
		return outputYAML(data)
	}

	return outputJSON(data)
}
