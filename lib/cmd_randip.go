package lib

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/pflag"
)

// CmdRandIPFlags are flags expected by CmdRandIP.
type CmdRandIPFlags struct {
	Help bool
	N    int
	IPv4 bool
	IPv6 bool
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
		&f.N,
		"count", "n", 1,
		"number of IPs to generate",
	)
	pflag.BoolVarP(
		&f.IPv4,
		"ipv4", "4", false,
		"generates ipv4 IPs",
	)
	pflag.BoolVarP(
		&f.IPv6,
		"ipv6", "6", false,
		"generates ipv6 IPs",
	)

}

func CmdRandIP(f CmdRandIPFlags, args []string, printHelp func()) error {
	if f.Help {
		printHelp()
		return nil
	}

	if f.IPv4 && f.IPv6 {
		return fmt.Errorf("only ipv4 or ipv6 allowed, but not both")
	} else if !f.IPv4 && !f.IPv6 {
		f.IPv4 = true
	}

	rand.Seed(time.Now().Unix())
	if f.IPv4 {
		RandIP4ListWrite(f.N)
	} else if f.IPv6 {
		RandIP6ListWrite(f.N)
	}
	return nil
}
