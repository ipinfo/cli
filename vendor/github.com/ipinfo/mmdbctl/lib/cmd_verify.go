package lib

import (
	"errors"
	"fmt"

	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/pflag"
)

// CmdVerifyFlags are flags expected by CmdVerify.
type CmdVerifyFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdVerify with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdVerifyFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)

}

func CmdVerify(f CmdVerifyFlags, args []string, printHelp func()) error {
	// help?
	if f.Help || (pflag.NArg() == 1 && pflag.NFlag() == 0) {
		printHelp()
		return nil
	}

	// validate input file.
	if len(args) == 0 {
		return errors.New("input mmdb file required as first argument")
	}

	// open tree.
	db, err := maxminddb.Open(args[0])
	if err != nil {
		return fmt.Errorf("couldn't open mmdb file: %w", err)
	}
	defer db.Close()

	// verify.
	err = db.Verify()
	if err != nil {
		fmt.Printf("invalid: %v\n", err)
	} else {
		fmt.Println("valid")
	}

	return nil
}
