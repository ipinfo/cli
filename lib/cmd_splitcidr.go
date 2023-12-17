package lib

import (
	"fmt"
	"net"
	"strconv"

	"github.com/ipinfo/cli/lib/iputil"
	"github.com/spf13/pflag"
)

type CmdSplitCIDRFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdSplitCIDR with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdSplitCIDRFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdSplitCIDR is the common core logic for the splitcidr command.
func CmdSplitCIDR(
	f CmdSplitCIDRFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	if len(args) != 2 {
		printHelp()
		return nil
	}

	cidrString := args[0]
	splitString := args[1]
	ip, _, err := net.ParseCIDR(cidrString)
	if err != nil {
		return err
	}

	split, err := strconv.Atoi(splitString)
	if err != nil {
		return nil
	}

	if ip.To4() != nil {
		ipsubnet, err := iputil.IPSubnetFromCidr(cidrString)
		if err != nil {
			return err
		}

		subs, err := ipsubnet.SplitCIDR(split)
		if err != nil {
			return err
		}

		for _, s := range subs {
			fmt.Println(s.ToCIDR())
		}
	} else {
		ipsubnet, err := iputil.IP6SubnetFromCidr(cidrString)
		if err != nil {
			return err
		}
		subs, err := ipsubnet.SplitCIDR(split)
		if err != nil {
			return err
		}
		for _, s := range subs {
			fmt.Println(s.ToCIDR())
		}
	}

	return nil
}
