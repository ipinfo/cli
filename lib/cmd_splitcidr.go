package lib

import (
	"fmt"
	"os"
	"strconv"

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
	// require args and/or stdin.
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if (len(args) == 0 || len(args) != 2) && !isStdin {
		printHelp()
		return nil
	}
	cidrString := args[0]
	splitString := args[1]
	cid, _ := CIDRToIPSubnet(cidrString)
	split, err := strconv.Atoi(splitString)
	if err != nil {
		return nil
	}
	bitshifts := split - int(cid.NetBitCnt)
	if bitshifts < 0 || bitshifts > 31 || int(cid.NetBitCnt)+bitshifts > 32 {
		return fmt.Errorf("wrong split string")
	}
	cids, _ := cid.SubnetBitShift(bitshifts)
	for _, cidr := range cids {
		fmt.Println(cidr.ToCIDR())
	}

	return nil
}
