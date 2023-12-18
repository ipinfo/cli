package mmdbwriter

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
)

// AliasedNetworkError is returned when inserting a aliased network into
// a Tree where DisableIPv4Aliasing in Options is false.
type AliasedNetworkError struct {
	// AliasedNetwork is the aliased network being inserted into.
	AliasedNetwork netip.Prefix
	// InsertedNetwork is the network being inserted into the Tree.
	InsertedNetwork netip.Prefix
}

func newAliasedNetworkError(netIP net.IP, curPrefixLen, recPrefixLen int) error {
	anErr := &AliasedNetworkError{}
	ip, ok := netip.AddrFromSlice(netIP)
	if !ok {
		return errors.Join(
			fmt.Errorf("creating netip.Addr from %s", netIP),
			anErr,
		)
	}
	var err error
	// We are using netip here despite using net.IP/net.IPNet internally as
	// it seems quite likely that we will switch to netip throughout.
	anErr.InsertedNetwork, err = ip.Prefix(recPrefixLen)
	if err != nil {
		return errors.Join(
			fmt.Errorf(
				"creating prefix from addr %s and prefix length %d: %w",
				ip,
				recPrefixLen,
				err,
			),
			anErr,
		)
	}

	anErr.AliasedNetwork, err = ip.Prefix(curPrefixLen)
	if err != nil {
		return errors.Join(
			fmt.Errorf(
				"creating prefix from addr %s and prefix length %d: %w",
				ip,
				curPrefixLen,
				err,
			),
			anErr,
		)
	}
	return anErr
}

func (r *AliasedNetworkError) Error() string {
	return fmt.Sprintf(
		"attempt to insert %s into %s, which is an aliased network",
		r.InsertedNetwork,
		r.AliasedNetwork,
	)
}

// ReservedNetworkError is returned when inserting a reserved network into
// a Tree where IncludeReservedNetworks in Options is false.
type ReservedNetworkError struct {
	// InsertedNetwork is the network being inserted into the Tree.
	InsertedNetwork netip.Prefix
	// ReservedNetwork is the reserved network being inserted into.
	ReservedNetwork netip.Prefix
}

var _ error = &ReservedNetworkError{}

func newReservedNetworkError(netIP net.IP, curPrefixLen, recPrefixLen int) error {
	rnErr := &ReservedNetworkError{}
	ip, ok := netip.AddrFromSlice(netIP)
	if !ok {
		return errors.Join(
			fmt.Errorf("creating netip.Addr from %s", netIP),
			rnErr,
		)
	}
	var err error
	// We are using netip here despite using net.IP/net.IPNet internally as
	// it seems quite likely that we will switch to netip throughout.
	rnErr.InsertedNetwork, err = ip.Prefix(recPrefixLen)
	if err != nil {
		return errors.Join(
			fmt.Errorf(
				"creating prefix from addr %s and prefix length %d: %w",
				ip,
				recPrefixLen,
				err,
			),
			rnErr,
		)
	}

	rnErr.ReservedNetwork, err = ip.Prefix(curPrefixLen)
	if err != nil {
		return errors.Join(
			fmt.Errorf(
				"creating prefix from addr %s and prefix length %d: %w",
				ip,
				curPrefixLen,
				err,
			),
			rnErr,
		)
	}
	return rnErr
}

func (r *ReservedNetworkError) Error() string {
	return fmt.Sprintf(
		"attempt to insert %s into %s, which is a reserved network",
		r.InsertedNetwork,
		r.ReservedNetwork,
	)
}
