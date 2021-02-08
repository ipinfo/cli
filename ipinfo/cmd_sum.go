package main

import (
	"github.com/spf13/pflag"
)

func printHelpSum() {
	// TODO
}

func cmdSum() error {
	var fTok string
	var fHelp bool
	var fPretty bool
	var fJSON bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVarP(&fPretty, "pretty", "p", true, "output pretty format. (default)")
	pflag.BoolVarP(&fJSON, "json", "j", true, "output JSON format.")
	pflag.Parse()

	if fHelp {
		printHelpSum()
		return nil
	}

	if err := prepareIpinfoClient(fTok); err != nil {
		return err
	}

	// TODO

	return nil
}
