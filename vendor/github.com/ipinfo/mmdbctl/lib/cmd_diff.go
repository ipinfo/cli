package lib

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"reflect"

	"github.com/oschwald/maxminddb-golang"
	"github.com/spf13/pflag"
)

// CmdDiffFlags are flags expected by CmdDiff.
type CmdDiffFlags struct {
	Help    bool
	Subnets bool
	Records bool
}

// Init initializes the common flags available to CmdDiff with sensible
// defaults.
//
// pflag.Parse() must be called to actually use the final flag values.
func (f *CmdDiffFlags) Init() {
	_h := "see description in --help"
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVarP(
		&f.Subnets,
		"subnets", "s", false,
		_h,
	)
	pflag.BoolVarP(
		&f.Records,
		"records", "r", false,
		_h,
	)
}

type cmdDiffRecord struct {
	oldr interface{}
	newr interface{}
}

func doDiff(
	newDb *maxminddb.Reader,
	newDbStr string,
	oldDb *maxminddb.Reader,
	oldDbStr string,
) (map[interface{}]*net.IPNet, map[interface{}]cmdDiffRecord, error) {
	modifiedSubnets := map[interface{}]*net.IPNet{}
	modifiedRecords := map[interface{}]cmdDiffRecord{}
	networksA := newDb.Networks(maxminddb.SkipAliasedNetworks)
	for networksA.Next() {
		var recordA interface{}
		var recordB interface{}

		subnetA, err := networksA.Network(&recordA)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to get record for subnet from %v: %w",
				newDbStr, err,
			)
		}

		subnetB, _, err := oldDb.LookupNetwork(subnetA.IP, &recordB)
		if err != nil {
			return nil, nil, fmt.Errorf(
				"failed to get record for IP %v from %v: %w",
				subnetA.IP, oldDbStr, err,
			)
		}

		// unequal subnets?
		if bytes.Compare(subnetA.IP, subnetB.IP) != 0 ||
			bytes.Compare(subnetA.Mask, subnetB.Mask) != 0 {
			modifiedSubnets[subnetA] = subnetB
			continue
		}

		// different data for same subnet?
		if !reflect.DeepEqual(recordA, recordB) {
			modifiedRecords[subnetA] = cmdDiffRecord{
				oldr: recordB,
				newr: recordA,
			}
		}
	}
	if networksA.Err() != nil {
		return nil, nil, fmt.Errorf(
			"failed traversing networks of %v: %w",
			newDbStr, networksA.Err(),
		)
	}

	return modifiedSubnets, modifiedRecords, nil
}

func CmdDiff(f CmdDiffFlags, args []string, printHelp func()) error {
	if f.Help || (pflag.NArg() == 1 && pflag.NFlag() == 0) {
		printHelp()
		return nil
	}

	// validate input files.
	if len(args) != 2 {
		return errors.New("two input mmdb file required as arguments")
	}

	// open old db.
	oldMmdb := args[0]
	oldDb, err := maxminddb.Open(oldMmdb)
	if err != nil {
		return fmt.Errorf("couldnt open %v: %w", oldMmdb, err)
	}
	defer oldDb.Close()

	// open new db.
	newMmdb := args[1]
	newDb, err := maxminddb.Open(newMmdb)
	if err != nil {
		return fmt.Errorf("couldnt open %v: %w", newMmdb, err)
	}
	defer newDb.Close()

	// confirm that they're of the same IP version.
	if newDb.Metadata.IPVersion != oldDb.Metadata.IPVersion {
		return fmt.Errorf(
			"IP versions differ between files: %v=%v and %v=%v",
			newMmdb, newDb.Metadata.IPVersion,
			oldMmdb, oldDb.Metadata.IPVersion,
		)
	}

	// collect set difference data.
	ambSn, ambRec, err := doDiff(newDb, newMmdb, oldDb, oldMmdb)
	if err != nil {
		return err
	}
	bmaSn, _, err := doDiff(oldDb, oldMmdb, newDb, newMmdb)
	if err != nil {
		return err
	}

	// print.
	if f.Subnets {
		if len(ambSn) > 0 || len(bmaSn) > 0 {
			fmt.Println("** SUBNETS **")
			for newSn, oldSn := range ambSn {
				fmt.Printf("%v -> %v\n", oldSn, newSn)
			}
			for newSn, oldSn := range bmaSn {
				fmt.Printf("%v -> %v\n", newSn, oldSn)
			}
		}
		fmt.Println(len(ambSn)+len(bmaSn), "subnet(s) modified.")
	}
	if f.Records {
		if f.Subnets {
			fmt.Println()
		}

		if len(ambRec) > 0 {
			fmt.Println("** RECORDS **")
			for sn, diffRecord := range ambRec {
				fmt.Println(sn)
				fmt.Printf("	-%v\n", diffRecord.oldr)
				fmt.Printf("	+%v\n", diffRecord.newr)
			}
		}
		fmt.Println(len(ambRec), "record(s) modified.")
	}
	if !f.Subnets && !f.Records {
		fmt.Println(len(ambSn)+len(bmaSn), "subnet(s) modified.")
		fmt.Println(len(ambRec), "record(s) modified.")
	}

	return nil
}
