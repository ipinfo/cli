package lib

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/spf13/pflag"
)

// CmdImportFlags are flags expected by CmdImport.
type CmdImportFlags struct {
	Help                bool
	In                  string
	Out                 string
	Csv                 bool
	Tsv                 bool
	Json                bool
	Fields              []string
	FieldsFromHdr       bool
	RangeMultiCol       bool
	JoinKeyCol          bool
	NoFields            bool
	NoNetwork           bool
	Ip                  int
	Size                int
	Merge               string
	IgnoreEmptyVals     bool
	DisallowReserved    bool
	Alias6to4           bool
	DisableMetadataPtrs bool
}

// Init initializes the common flags available to CmdImport with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdImportFlags) Init() {
	_h := "see description in --help"
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.StringVarP(
		&f.In,
		"in", "i", "",
		_h,
	)
	pflag.StringVarP(
		&f.Out,
		"out", "o", "",
		_h,
	)
	pflag.BoolVarP(
		&f.Csv,
		"csv", "c", false,
		_h,
	)
	pflag.BoolVarP(
		&f.Tsv,
		"tsv", "t", false,
		_h,
	)
	pflag.BoolVarP(
		&f.Json,
		"json", "j", false,
		_h,
	)
	pflag.StringSliceVarP(
		&f.Fields,
		"fields", "f", nil,
		_h,
	)
	pflag.BoolVar(
		&f.FieldsFromHdr,
		"fields-from-header", false,
		_h,
	)
	pflag.BoolVar(
		&f.RangeMultiCol,
		"range-multicol", false,
		_h,
	)
	pflag.BoolVar(
		&f.JoinKeyCol,
		"joinkey-col", false,
		_h,
	)
	pflag.BoolVar(
		&f.NoFields,
		"no-fields", false,
		_h,
	)
	pflag.BoolVar(
		&f.NoNetwork,
		"no-network", false,
		_h,
	)
	pflag.IntVar(
		&f.Ip,
		"ip", 6,
		_h,
	)
	pflag.IntVarP(
		&f.Size,
		"size", "s", 32,
		_h,
	)
	pflag.StringVarP(
		&f.Merge,
		"merge", "m", "none",
		_h,
	)
	pflag.BoolVar(
		&f.IgnoreEmptyVals,
		"ignore-empty-values", false,
		_h,
	)
	pflag.BoolVar(
		&f.DisallowReserved,
		"disallow-reserved", false,
		_h,
	)
	pflag.BoolVar(
		&f.Alias6to4,
		"alias-6to4", false,
		_h,
	)
	pflag.BoolVar(
		&f.DisableMetadataPtrs,
		"disable-metadata-pointers", true,
		_h,
	)

}

func CmdImport(f CmdImportFlags, args []string, printHelp func()) error {
	// help?
	if f.Help || (pflag.NArg() == 1 && pflag.NFlag() == 0) {
		printHelp()
		return nil
	}

	// optional input as 1st and output as 2nd argument.
	if len(args) >= 2 {
		f.In = args[0]
		f.Out = args[1]
	}

	// validate IP version.
	if f.Ip != 4 && f.Ip != 6 {
		return errors.New("ip version must be \"4\" or \"6\"")
	}

	// validate record size.
	if f.Size != 24 && f.Size != 28 && f.Size != 32 {
		return errors.New("record size must be 24, 28 or 32")
	}

	// validate merge strategy.
	var mergeStrategy inserter.FuncGenerator
	if f.Merge == "none" {
		mergeStrategy = inserter.ReplaceWith
	} else if f.Merge == "toplevel" {
		mergeStrategy = inserter.TopLevelMergeWith
	} else if f.Merge == "recurse" {
		mergeStrategy = inserter.DeepMergeWith
	} else {
		return errors.New("merge strategy must be \"none\", \"toplevel\" or \"recurse\"")
	}

	// figure out file type.
	var delim rune
	if !f.Csv && !f.Tsv && !f.Json {
		if strings.HasSuffix(f.In, ".csv") {
			delim = ','
		} else if strings.HasSuffix(f.In, ".tsv") {
			delim = '\t'
		} else if strings.HasSuffix(f.In, ".json") {
			delim = '-'
		} else {
			return errors.New("input file type unknown")
		}
	} else {
		if f.Csv && f.Tsv || f.Csv && f.Json || f.Tsv && f.Json {
			return errors.New("multiple input file types specified")
		} else if f.Csv {
			delim = ','
		} else if f.Tsv {
			delim = '\t'
		} else {
			delim = '-'
		}
	}

	// figure out fields.
	fieldSrcCnt := 0
	if f.Fields != nil && len(f.Fields) > 0 {
		fieldSrcCnt += 1
	}
	if f.FieldsFromHdr {
		fieldSrcCnt += 1
	}
	if f.NoFields {
		fieldSrcCnt += 1
	}
	if fieldSrcCnt > 1 {
		return errors.New("conflicting field sources specified.")
	}
	if f.NoFields {
		f.Fields = []string{}
		f.NoNetwork = false
	} else if !f.FieldsFromHdr && (f.Fields == nil || len(f.Fields) == 0) {
		f.FieldsFromHdr = true
	}

	if f.JoinKeyCol {
		f.RangeMultiCol = true
	}

	// prepare output file.
	var outFile *os.File
	if f.Out == "" {
		outFile = os.Stdout
	} else {
		var err error
		outFile, err = os.Create(f.Out)
		if err != nil {
			return fmt.Errorf("could not create %v: %w", f.Out, err)
		}
		defer outFile.Close()
	}

	// init tree.
	dbdesc := "ipinfo " + filepath.Base(f.Out)
	tree, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: dbdesc,
			Description: map[string]string{
				"en": dbdesc,
			},
			Languages:               []string{"en"},
			DisableIPv4Aliasing:     !f.Alias6to4,
			IncludeReservedNetworks: !f.DisallowReserved,
			IPVersion:               f.Ip,
			RecordSize:              f.Size,
			DisableMetadataPointers: f.DisableMetadataPtrs,
			Inserter:                mergeStrategy,
		},
	)
	if err != nil {
		return fmt.Errorf("could not create tree: %w", err)
	}

	// prepare input file.
	var inFile *os.File
	if f.In == "" || f.In == "-" {
		inFile = os.Stdin
	} else {
		var err error
		inFile, err = os.Open(f.In)
		if err != nil {
			return fmt.Errorf("invalid input file %v: %w", f.In, err)
		}
		defer inFile.Close()
	}

	inFileBuffered := bufio.NewReaderSize(inFile, 65536)

	entrycnt := 0
	if delim == ',' || delim == '\t' {
		var rdr reader
		if delim == ',' {
			csvrdr := csv.NewReader(inFileBuffered)
			csvrdr.Comma = delim
			csvrdr.LazyQuotes = true

			rdr = csvrdr
		} else {
			tsvrdr := NewTsvReader(inFileBuffered)

			rdr = tsvrdr
		}

		// read from input, scanning & parsing each line according to delim,
		// then insert that into the tree.
		dataColStart := 1
		hdrSeen := false
		for {
			parts, err := rdr.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("input scanning failed: %w", err)
			}

			// on header line?
			if !hdrSeen {
				hdrSeen = true

				// check if the header has a multi-column range.
				if len(parts) > 1 && parts[0] == "start_ip" && parts[1] == "end_ip" {
					f.RangeMultiCol = true

					// maybe we also have a join key?
					if len(parts) > 2 && parts[2] == "join_key" {
						f.JoinKeyCol = true
					}
				}

				if f.RangeMultiCol {
					if f.JoinKeyCol {
						dataColStart = 3
					} else {
						dataColStart = 2
					}
				}

				// need to get fields from hdr?
				if f.FieldsFromHdr {
					// skip all non-data columns.
					f.Fields = parts[dataColStart:]
				}

				// insert empty values for all fields in 0.0.0.0/0 if requested.
				if f.IgnoreEmptyVals {
					_, network, _ := net.ParseCIDR("0.0.0.0/0")
					record := mmdbtype.Map{}
					for _, field := range f.Fields {
						record[mmdbtype.String(field)] = mmdbtype.String("")
					}
					if err := tree.Insert(network, record); err != nil {
						return errors.New(
							"couldn't insert empty values to 0.0.0.0/0",
						)
					}
				}

				// should we skip this first line now?
				if f.FieldsFromHdr {
					continue
				}
			}

			networkStr := parts[0]

			// convert 2 IPs into IP range?
			if f.RangeMultiCol {
				networkStr = parts[0] + "-" + parts[1]
			}

			// add network part to single-IP network if it's missing.
			isNetworkRange := strings.Contains(networkStr, "-")
			if !isNetworkRange && !strings.Contains(networkStr, "/") {
				if f.Ip == 6 && strings.Contains(networkStr, ":") {
					networkStr += "/128"
				} else {
					networkStr += "/32"
				}
			}

			// prep record.
			record := mmdbtype.Map{}
			if !f.NoNetwork {
				record["network"] = mmdbtype.String(networkStr)
			}
			for i, field := range f.Fields {
				record[mmdbtype.String(field)] = mmdbtype.String(parts[i+dataColStart])
			}

			// range insertion or cidr insertion?
			if isNetworkRange {
				networkStrParts := strings.Split(networkStr, "-")
				startIp := net.ParseIP(networkStrParts[0])
				endIp := net.ParseIP(networkStrParts[1])
				if err := tree.InsertRange(startIp, endIp, record); err != nil {
					fmt.Fprintf(
						os.Stderr, "warn: couldn't insert line '%v'\n",
						strings.Join(parts, string(delim)),
					)
				}
			} else {
				_, network, err := net.ParseCIDR(networkStr)
				if err != nil {
					return fmt.Errorf(
						"couldn't parse cidr \"%v\": %w",
						networkStr, err,
					)
				}
				if err := tree.Insert(network, record); err != nil {
					fmt.Fprintf(
						os.Stderr, "warn: couldn't insert line '%v'\n",
						strings.Join(parts, string(delim)),
					)
				}
			}

			entrycnt += 1
		}
	} else if delim == '-' {
		dataStream := json.NewDecoder(inFileBuffered)
		for {
			// Decode one JSON document.
			var row interface{}
			err := dataStream.Decode(&row)

			if err != nil {
				// io.EOF is expected at end of stream.
				if err != io.EOF {
					return fmt.Errorf("error in io.EOF: %w", err)
				}
				break
			}
			mResult := row.(map[string]interface{})

			// insert empty values for all fields in 0.0.0.0/0 if requested.
			if f.IgnoreEmptyVals {
				_, network, _ := net.ParseCIDR("0.0.0.0/0")
				record := mmdbtype.Map{}
				for _, field := range f.Fields {
					record[mmdbtype.String(field)] = mmdbtype.String("")
				}
				if err := tree.Insert(network, record); err != nil {
					return errors.New(
						"couldn't insert empty values to 0.0.0.0/0",
					)
				}
			}

			// convert 2 IPs into IP range?
			var networkStr string
			if val, ok := mResult["start_ip"].(string); ok {
				networkStr = val + "-" + mResult["end_ip"].(string)
				delete(mResult, "start_ip")
				delete(mResult, "end_ip")
				if _, ok := mResult["join_key"].(string); ok {
					delete(mResult, "join_key")
				}
			} else if val, ok := mResult["range"].(string); ok {
				networkStr = val
				delete(mResult, "range")
			} else {
				return errors.New(
					"couldn't get ip or range from the record",
				)
			}

			// add network part to single-IP network if it's missing.
			isNetworkRange := strings.Contains(networkStr, "-")
			if !isNetworkRange && !strings.Contains(networkStr, "/") {
				if f.Ip == 6 && strings.Contains(networkStr, ":") {
					networkStr += "/128"
				} else {
					networkStr += "/32"
				}
			}

			// prep record.
			record := mmdbtype.Map{}
			if !f.NoNetwork {
				record["network"] = mmdbtype.String(networkStr)
			}

			mResultStr := mapInterfaceToStr(mResult)
			for k, v := range mResultStr {
				record[mmdbtype.String(k)] = mmdbtype.String(v)
			}

			// range insertion or cidr insertion?
			if isNetworkRange {
				networkStrParts := strings.Split(networkStr, "-")
				startIp := net.ParseIP(networkStrParts[0])
				endIp := net.ParseIP(networkStrParts[1])
				if err := tree.InsertRange(startIp, endIp, record); err != nil {
					fmt.Fprintf(
						os.Stderr, "warn: couldn't insert '%v'\n",
						mResult,
					)
				}
			} else {
				_, network, err := net.ParseCIDR(networkStr)
				if err != nil {
					return fmt.Errorf(
						"couldn't parse cidr \"%v\": %w",
						networkStr, err,
					)
				}
				if err := tree.Insert(network, record); err != nil {
					fmt.Fprintf(
						os.Stderr, "warn: couldn't insert '%v'\n",
						mResult,
					)
				}
			}

			entrycnt += 1
		}
	}

	if entrycnt == 0 {
		return errors.New("nothing to import")
	}

	// write out mmdb file.
	fmt.Fprintf(os.Stderr, "writing to %s (%v entries)\n", f.Out, entrycnt)
	if _, err := tree.WriteTo(outFile); err != nil {
		return fmt.Errorf("writing out to tree failed: %w", err)
	}

	return nil
}
