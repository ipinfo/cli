package lib

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// CmdCIDR2RangeFlags are flags expected by CmdCIDR2Range.
type CmdCIDR2RangeFlags struct {
	Help bool
}

// Init initializes the common flags available to CmdCIDR2Range with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdCIDR2RangeFlags) Init() {
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
}

// CmdCIDR2Range is the common core logic for the cidr2range command.
func CmdCIDR2Range(
	f CmdCIDR2RangeFlags,
	args []string,
	printHelp func(),
) error {
	if f.Help {
		printHelp()
		return nil
	}

	// require args.
	stat, _ := os.Stdin.Stat()
	isStdin := (stat.Mode() & os.ModeCharDevice) == 0
	if len(args) == 0 && !isStdin {
		printHelp()
		return nil
	}

	// actual scanner.
	scanrdr := func(r io.Reader) {
		var rem string
		var hitEOF bool

		buf := bufio.NewReader(r)
		for {
			if hitEOF {
				return
			}

			d, err := buf.ReadString('\n')
			if err == io.EOF {
				if len(d) == 0 {
					return
				}

				// do one more loop on remaining content.
				hitEOF = true
			} else if err != nil {
				// TODO print error but have a `-q` flag to be quiet.
				return
			}

			sepIdx := strings.IndexAny(d, ",\n")
			if sepIdx == -1 {
				// only possible if EOF & input doesn't end with newline.
				sepIdx = len(d)
				rem = "\n"
			} else {
				rem = d[sepIdx:]
			}

			cidrStr := d[:sepIdx]
			if strings.IndexByte(cidrStr, ':') == -1 {
				if r, err := IPRangeStrFromCIDR(cidrStr); err == nil {
					fmt.Printf("%s%s", r.String(), rem)
				} else {
					goto noip
				}
			} else {
				if r, err := IP6RangeStrFromCIDR(cidrStr); err == nil {
					fmt.Printf("%s%s", r.String(), rem)
				} else {
					goto noip
				}
			}

			continue

		noip:
			fmt.Printf("%s", d)
			if sepIdx == len(d) {
				fmt.Println()
			}
		}
	}

	// scan stdin first.
	if isStdin {
		scanrdr(os.Stdin)
	}

	// scan all args.
	for _, arg := range args {
		f, err := os.Open(arg)
		if err != nil {
			// is it a CIDR?
			if r, err := IPRangeStrFromCIDR(arg); err == nil {
				fmt.Println(r.String())
				continue
			}

			// invalid file arg.
			return err
		}

		scanrdr(f)
	}

	return nil
}
