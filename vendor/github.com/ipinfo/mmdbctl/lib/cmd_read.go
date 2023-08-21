package lib

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib"
	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/pflag"
)

var predictReadFmts = []string{
	"json",
	"json-compact",
	"json-pretty",
	"tsv",
	"csv",
}

// CmdReadFlags are flags expected by CmdRead.
type CmdReadFlags struct {
	Help    bool
	NoColor bool
	Format  string
}

// Init initializes the common flags available to CmdRead with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdReadFlags) Init() {
	_h := "see description in --help"
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVar(
		&f.NoColor,
		"nocolor", false,
		_h,
	)
	pflag.StringVarP(
		&f.Format,
		"format", "f", "json",
		_h,
	)
}

func CmdRead(f CmdReadFlags, args []string, printHelp func()) error {
	if f.NoColor {
		color.NoColor = true
	}

	// help?
	if f.Help || (pflag.NArg() == 1 && pflag.NFlag() == 0) {
		printHelp()
		return nil
	}

	// validate format.
	if f.Format == "json" {
		f.Format = "json-compact"
	}
	validFormat := false
	for _, format := range predictReadFmts {
		if f.Format == format {
			validFormat = true
			break
		}
	}
	if !validFormat {
		return fmt.Errorf("format must be one of %v", predictReadFmts)
	}

	// last arg must be mmdb file; open it.
	mmdbFileArg := args[len(args)-1]
	db, err := maxminddb.Open(mmdbFileArg)
	if err != nil {
		return fmt.Errorf("couldn't open mmdb file %v: %w", mmdbFileArg, err)
	}
	defer db.Close()

	// get IP list.
	ips, err := lib.IPListFromAllSrcs(args[:len(args)-1])
	if err != nil {
		return fmt.Errorf("couldn't get IP list: %w", err)
	}

	requiresHdr := f.Format == "csv" || f.Format == "tsv"
	hdrWritten := false
	var wr writer
	if f.Format == "csv" {
		csvwr := csv.NewWriter(os.Stdout)
		wr = csvwr
	} else if f.Format == "tsv" {
		tsvwr := NewTsvWriter(os.Stdout)
		wr = tsvwr
	}
	for _, ip := range ips {
		record := make(map[string]interface{})
		if err := db.Lookup(ip, &record); err != nil || len(record) == 0 {
			if !requiresHdr {
				fmt.Fprintf(os.Stderr,
					"err: couldn't get data for %s\n",
					ip.String(),
				)
			}
			continue
		}
		recordStr := mapInterfaceToStr(record)

		if !hdrWritten {
			hdrWritten = true

			if requiresHdr {
				hdr := append([]string{"ip"}, sortedMapKeys(recordStr)...)
				if err := wr.Write(hdr); err != nil {
					return fmt.Errorf(
						"failed to write header %v: %w",
						hdr, err,
					)
				}
			}
		}

		if f.Format == "json-compact" {
			b, err := json.Marshal(record)
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"err: couldn't print data for %s\n",
					ip.String(),
				)
				continue
			}
			fmt.Printf("%s\n", b)
		} else if f.Format == "json-pretty" {
			b, err := json.MarshalIndent(record, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr,
					"err: couldn't print data for %s\n",
					ip.String(),
				)
				continue
			}
			fmt.Printf("%s\n", b)
		} else { // if fFormat == "csv" || fFormat == "tsv"
			line := append([]string{ip.String()}, sortedMapValsByKeys(recordStr)...)
			if err := wr.Write(line); err != nil {
				return fmt.Errorf("failed to write line %v: %w", line, err)
			}
		}
	}
	if wr != nil {
		wr.Flush()
		if err := wr.Error(); err != nil {
			return fmt.Errorf("writer had failure: %w", err)
		}
	}

	return nil
}
