package lib

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/pflag"
)

// CmdExportFlags are flags expected by CmdExport.
type CmdExportFlags struct {
	Help   bool
	NoHdr  bool
	Format string
	Out    string
}

// Init initializes the common flags available to CmdExport with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdExportFlags) Init() {
	_h := "see description in --help"
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVar(
		&f.NoHdr,
		"no-header", false,
		_h,
	)
	pflag.StringVarP(
		&f.Format,
		"format", "f", "",
		_h,
	)
	pflag.StringVarP(
		&f.Out,
		"out", "o", "",
		_h,
	)
}

func CmdExport(f CmdExportFlags, args []string, printHelp func()) error {
	// help?
	if f.Help || (pflag.NArg() == 1 && pflag.NFlag() == 0) {
		printHelp()
		return nil
	}

	// validate input file.
	if len(args) == 0 {
		return errors.New("input mmdb file required as first argument")
	}

	// prepare output file.
	var outFile *os.File
	if f.Out == "" && len(args) < 2 {
		outFile = os.Stdout
	} else {
		// either flag or argument is defined.
		if f.Out == "" {
			f.Out = args[1]
		}

		var err error
		outFile, err = os.Create(f.Out)
		if err != nil {
			return fmt.Errorf("could not create %v: %w", f.Out, err)
		}
		defer outFile.Close()
	}

	// validate format.
	if f.Format == "" {
		if strings.HasSuffix(f.Out, ".csv") {
			f.Format = "csv"
		} else if strings.HasSuffix(f.Out, ".tsv") {
			f.Format = "tsv"
		} else if strings.HasSuffix(f.Out, ".json") {
			f.Format = "json"
		} else {
			f.Format = "csv"
		}
	}
	if f.Format != "csv" && f.Format != "tsv" && f.Format != "json" {
		return errors.New("format must be \"csv\" or \"tsv\" or \"json\"")
	}

	// open tree.
	db, err := maxminddb.Open(args[0])
	if err != nil {
		return fmt.Errorf("couldn't open mmdb file: %w", err)
	}
	defer db.Close()

	if f.Format == "tsv" || f.Format == "csv" {
		// export.
		hdrWritten := false
		var wr writer
		if f.Format == "csv" {
			csvwr := csv.NewWriter(outFile)
			wr = csvwr
		} else {
			tsvwr := NewTsvWriter(outFile)
			wr = tsvwr
		}
		record := make(map[string]interface{})
		networks := db.Networks(maxminddb.SkipAliasedNetworks)
		for networks.Next() {
			subnet, err := networks.Network(&record)
			if err != nil {
				return fmt.Errorf("failed to get record for next subnet: %w", err)
			}

			recordStr := mapInterfaceToStr(record)
			if !hdrWritten {
				hdrWritten = true

				if !f.NoHdr {
					hdr := append([]string{"range"}, sortedMapKeys(recordStr)...)
					if err := wr.Write(hdr); err != nil {
						return fmt.Errorf(
							"failed to write header %v: %w",
							hdr, err,
						)
					}
				}
			}

			line := append(
				[]string{subnet.String()},
				sortedMapValsByKeys(recordStr)...,
			)
			if err := wr.Write(line); err != nil {
				return fmt.Errorf("failed to write line %v: %w", line, err)
			}
		}
		wr.Flush()
		if err := wr.Error(); err != nil {
			return fmt.Errorf("writer had failure: %w", err)
		}
		if err := networks.Err(); err != nil {
			return fmt.Errorf("failed networks traversal: %w", err)
		}
	} else if f.Format == "json" {
		networks := db.Networks(maxminddb.SkipAliasedNetworks)
		enc := json.NewEncoder(outFile)
		for networks.Next() {
			record := make(map[string]interface{})

			subnet, err := networks.Network(&record)
			if err != nil {
				return fmt.Errorf("failed to get record for next subnet: %w", err)
			}
			record["range"] = subnet.String()
			enc.Encode(record)
		}
		if err := networks.Err(); err != nil {
			return fmt.Errorf("failed networks traversal: %w", err)
		}
	}
	return nil
}
