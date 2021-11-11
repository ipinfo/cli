package lib

import (
	"math/rand"
	"time"

	"github.com/spf13/pflag"
)

// CmdRandIPFlags are flags expected by CmdRandIP.
type CmdRandIPFlags struct {
	Help bool
	n    int
	Type string
}

// Init initializes the common flags available to CmdRandIP with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdRandIPFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.IntVarP(
		&f.n,
		"count", "n", 1,
		"number of IPs to generate",
	)
	pflag.StringVarP(
		&f.Type,
		"type", "t", "ipv4",
		"ipv4/ipv6",
	)

}

func CmdRandIP(f CmdRandIPFlags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	rand.Seed(time.Now().Unix())
	if f.Type == "ipv4" || f.Type == "IPV4" {
		RandIP4Write(f.n)
	} else {
		printHelp()
	}
	return nil
}
