package lib

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/pflag"
)

// CmdRandIPFlags are flags expected by CmdRandIP.
type CmdRandIPFlags struct {
	Help       bool
	N          int
	IPv4       bool
	IPv6       bool
	ExcludeRes bool
	StartIP    string
	EndIP      string
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
		"num", "n", 1,
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
	pflag.BoolVarP(
		&f.ExcludeRes,
		"exclude-reserved", "x", false,
		"exclude reserved/bogon IPs.",
	)
	pflag.StringVarP(
		&f.StartIP,
		"start", "s", "",
		"starting range of IPs",
	)
	pflag.StringVarP(
		&f.EndIP,
		"end", "e", "",
		"ending range of IPs",
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
		if f.StartIP == "" {
			f.StartIP = "0.0.0.0"
		}
		if f.EndIP == "" {
			f.EndIP = "255.255.255.255"
		}
		if err := RandIP4RangeListWrite(f.StartIP, f.EndIP, f.N, f.ExcludeRes); err != nil {
			return err
		}
	} else if f.IPv6 {
		if f.StartIP == "" {
			f.StartIP = "::"
		}
		if f.EndIP == "" {
			f.EndIP = "ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"
		}
		if err := RandIP6RangeListWrite(f.StartIP, f.EndIP, f.N, f.ExcludeRes); err != nil {
			return err
		}
	}
	return nil
}
