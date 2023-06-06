package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

var completionsMmdbCtl = &complete.Command{
	Sub: map[string]*complete.Command{
		"read":     completionsRead,
		"import":   completionsImport,
		"export":   completionsExport,
		"diff":     completionsDiff,
		"metadata": completionsMetadata,
		"verify":   completionsVerify,
	},
	Flags: map[string]complete.Predictor{
		"-h":     predict.Nothing,
		"--help": predict.Nothing,
	},
}

func printHelpMmdbCtl() {

	fmt.Printf(
		`Usage: %s mmdbctl <cmd> [<opts>] [<args>]

Commands:
  read        read data for IPs in an mmdb file.
  import      import data in non-mmdb format into mmdb.
  export      export data from mmdb format into non-mmdb format.
  diff        see the difference between two mmdb files.
  metadata    print metadata from the mmdb file.
  verify      check that the mmdb file is not corrupted or invalid.
  completion  install or output shell auto-completion script.

Options:
  General:
    --nocolor
      disable colored output.
    --help, -h
      show help.
`, progBase)
}

func mmdbctlHelp() (err error) {
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVar(&fNoColor, "nocolor", false, "disable colored output.")
	pflag.Parse()

	if fNoColor {
		color.NoColor = true
	}

	if fHelp {
		printHelpDefault()
		return nil
	}

	// currently we do nothing by default.
	printHelpMmdbCtl()
	return nil
}

func cmdMmdbCtl() error {
	var err error
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	switch {
	case cmd == "read":
		err = cmdRead()
	case cmd == "import":
		err = cmdImport()
	case cmd == "export":
		err = cmdExport()
	case cmd == "diff":
		err = cmdDiff()
	case cmd == "verify":
		err = cmdVerify()
	case cmd == "metadata":
		err = cmdMetadata()
	default:
		err = mmdbctlHelp()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %v\n", err)
	}

	return nil
}
