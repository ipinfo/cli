// Package mmdbwriter provides the tools to create and write MaxMind DB
// files.
package mmdbwriter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/maxmind/mmdbwriter/inserter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
	"github.com/oschwald/maxminddb-golang"
	"inet.af/netaddr"
)

var (
	metadataStartMarker  = []byte("\xAB\xCD\xEFMaxMind.com")
	dataSectionSeparator = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)

// Options holds configuration parameters for the writer.
type Options struct {
	// BuildEpoch is the database build timestamp as a Unix epoch value. It
	// defaults to the epoch of when New was called.
	BuildEpoch int64

	// DatabaseType is a string that indicates the structure of each data record
	// associated with an IP address. The actual definition of these structures
	// is left up to the database creator.
	DatabaseType string

	// Description is a map where the key is a language code and the value is
	// the description of the database in that language.
	Description map[string]string

	// DisableIPv4Aliasing will disable the IPv4 aliasing in IPv6 trees. This
	// aliasing maps some IPv6 networks to the IPv4 network, e.g.,
	// ::ffff:0:0/96.
	DisableIPv4Aliasing bool

	// IncludeReservedNetworks will allow reserved networks to be added to the
	// database.
	//
	// If this is false, any attempt to insert into these networks will result
	// in an error and inserting a network that contains a reserved network will
	// result in the reserved portion of the network being excluded. Reserved
	// networks that are globally routable to an individual device, such as
	// Teredo, may still be added.
	IncludeReservedNetworks bool

	// IPVersion indicates whether an IPv4 or IPv6 database should be built. An
	// IPv6 database supports both IPv4 and IPv6 lookups. The default value is
	// "6" for IPv6.
	IPVersion int

	// Languages is a slice of strings, each of which is a locale code. A given
	// record may contain data items that have been localized to some or all of
	// these locales. Records should not contain localized data for locales not
	// included in this slice.
	Languages []string

	// RecordSize indicates the number of bits in a record in the search tree.
	// The supported values are 24, 28, and 32. A smaller size will result in a
	// smaller database, but it will limit the maximum size of the database.
	// The default is 28.
	RecordSize int

	// DisableMetadataPointers prevents the use of pointers in the metadata
	// section of the database. This option exists to avoid bugs in reader
	// implementations that do not correctly handle metadata pointers. Its
	// use should primarily be limited to existing database types.
	DisableMetadataPointers bool

	// Inserter is the insert function used when calling `Insert`. It defaults
	// to `inserter.ReplaceWith`, which replaces any conflicting old value
	// entirely with the new.
	Inserter inserter.FuncGenerator
}

// Tree represents an MaxMind DB search tree.
type Tree struct {
	buildEpoch              int64
	databaseType            string
	dataMap                 *dataMap
	description             map[string]string
	disableMetadataPointers bool
	ipVersion               int
	languages               []string
	recordSize              int
	root                    *node
	treeDepth               int
	// This is set when the tree is finalized
	nodeCount       int
	inserterFuncGen inserter.FuncGenerator
}

// New creates a new Tree.
func New(opts Options) (*Tree, error) {
	tree := &Tree{
		buildEpoch:              time.Now().Unix(),
		dataMap:                 newDataMap(),
		databaseType:            opts.DatabaseType,
		description:             map[string]string{},
		disableMetadataPointers: opts.DisableMetadataPointers,
		ipVersion:               6,
		recordSize:              28,
		root:                    &node{},
		inserterFuncGen:         inserter.ReplaceWith,
	}

	if opts.BuildEpoch != 0 {
		tree.buildEpoch = opts.BuildEpoch
	}

	if opts.Description != nil {
		tree.description = opts.Description
	}

	if opts.IPVersion != 0 {
		tree.ipVersion = opts.IPVersion
	}

	if opts.Languages != nil {
		tree.languages = opts.Languages
	}

	if opts.RecordSize != 0 {
		tree.recordSize = opts.RecordSize
	}

	if opts.Inserter != nil {
		tree.inserterFuncGen = opts.Inserter
	}

	switch tree.ipVersion {
	case 6:
		tree.treeDepth = 128
	case 4:
		tree.treeDepth = 32
	default:
		return nil, fmt.Errorf("unsupported IPVersion: %d", tree.ipVersion)
	}

	if tree.ipVersion == 6 && !opts.DisableIPv4Aliasing {
		if err := tree.insertIPv4Aliases(); err != nil {
			return nil, err
		}
	}

	if !opts.IncludeReservedNetworks {
		err := tree.insertReservedNetworks()
		if err != nil {
			return nil, err
		}
	}

	return tree, nil
}

// Load an existing database into the writer.
func Load(path string, opts Options) (*Tree, error) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	metadata := db.Metadata
	if opts.DatabaseType == "" {
		opts.DatabaseType = metadata.DatabaseType
	}

	if opts.Description == nil {
		opts.Description = metadata.Description
	}

	if opts.IPVersion == 0 {
		opts.IPVersion = int(metadata.IPVersion)
	}

	if opts.Languages == nil {
		opts.Languages = metadata.Languages
	}

	if opts.RecordSize == 0 {
		opts.RecordSize = int(metadata.RecordSize)
	}

	tree, err := New(opts)
	if err != nil {
		return nil, err
	}

	dser := newDeserializer()

	var networkOpts []maxminddb.NetworksOption
	if opts.IPVersion == 6 && !opts.DisableIPv4Aliasing {
		networkOpts = append(networkOpts, maxminddb.SkipAliasedNetworks)
	}

	networks := db.Networks(networkOpts...)
	for networks.Next() {
		var network *net.IPNet

		dser.clear()
		network, err = networks.Network(dser)
		if err != nil {
			return nil, err
		}

		err = tree.Insert(network, dser.rv)
		if err != nil {
			return nil, err
		}
	}
	if err := networks.Err(); err != nil {
		return nil, err
	}
	return tree, nil
}

// Insert a data value into the tree using the Tree's inserter function
// (defaults to inserter.ReplaceWith).
//
// This is not safe to call from multiple threads.
func (t *Tree) Insert(network *net.IPNet, value mmdbtype.DataType) error {
	return t.InsertFunc(network, t.inserterFuncGen(value))
}

// InsertFunc will insert the output of the function passed to it. The argument
// passed to the function is the existing value in the record. The inserter
// function should return the mmdbtype.DataType to be inserted. In both cases,
// a nil value means an empty record.
//
// You must never modify the argument passed to the function as the value may
// be shared with other records. If you want a copy of the mmdbtype.DataType to modify,
// call the Copy method on it, which will make a deep copy. This isn't done
// automatically before calling the function as not all functions will require
// the record to be copied and there is a non-trivial performance impact.
//
// The function will be called multiple times per insert when the network
// has multiple preexisting records associated with it.
//
// This is not safe to call from multiple threads.
func (t *Tree) InsertFunc(
	network *net.IPNet,
	inserterFunc inserter.Func,
) error {
	return t.insert(network, recordTypeData, inserterFunc, nil)
}

func (t *Tree) insert(
	network *net.IPNet,
	recordType recordType,
	inserterFunc inserter.Func,
	node *node,
) error {
	// We set this to 0 so that the tree must be finalized again.
	t.nodeCount = 0

	prefixLen, _ := network.Mask.Size()

	ip := network.IP
	if t.treeDepth == 128 && len(ip) == 4 {
		ip = ipV4ToV6(ip)
		prefixLen += 96
	}

	return t.root.insert(
		insertRecord{
			ip:           ip,
			prefixLen:    prefixLen,
			recordType:   recordType,
			inserter:     inserterFunc,
			insertedNode: node,

			dataMap: t.dataMap,
		},
		0,
	)
}

// InsertRange is the same as Insert, except it will insert all subnets within
// the range of IPs specified by `[start,end]`.
func (t *Tree) InsertRange(
	start net.IP,
	end net.IP,
	value mmdbtype.DataType,
) error {
	return t.InsertRangeFunc(start, end, t.inserterFuncGen(value))
}

// InsertRangeFunc is the same as InsertFunc, except it will insert all subnets
// within the range of IPs specified by `[start,end]`.
func (t *Tree) InsertRangeFunc(
	start net.IP,
	end net.IP,
	inserterFunc inserter.Func,
) error {
	return t.insertRange(start, end, recordTypeData, inserterFunc, nil)
}

func (t *Tree) insertRange(
	start net.IP,
	end net.IP,
	recordType recordType,
	inserterFunc inserter.Func,
	node *node,
) error {
	_start, ok := netaddr.FromStdIP(start)
	if !ok {
		return errors.New("start IP is invalid")
	}
	_end, ok := netaddr.FromStdIP(end)
	if !ok {
		return errors.New("end IP is invalid")
	}

	r := netaddr.IPRangeFrom(_start, _end)
	if !r.IsValid() {
		return errors.New("start & end IPs did not give valid range")
	}
	subnets := r.Prefixes()
	for _, subnet := range subnets {
		if err := t.insert(subnet.IPNet(), recordType, inserterFunc, node); err != nil {
			return err
		}
	}

	return nil
}

func (t *Tree) insertStringNetwork(
	network string,
	recordType recordType,
	inserterFunc inserter.Func,
	node *node,
) error {
	_, ipnet, err := net.ParseCIDR(network)
	if err != nil {
		return fmt.Errorf("parsing network (%s): %w", network, err)
	}
	return t.insert(ipnet, recordType, inserterFunc, node)
}

var ipv4AliasNetworks = []string{
	"::ffff:0:0/96",
	"2001::/32",
	"2002::/16",
}

func (t *Tree) insertIPv4Aliases() error {
	_, ipv4Root, err := net.ParseCIDR("::/96")
	if err != nil {
		return fmt.Errorf("parsing IPv4 root: %w", err)
	}

	ipv4RootNode := &node{}

	// Make ::/96, the IPv4 root, a fixed node.
	err = t.insert(ipv4Root, recordTypeFixedNode, nil, ipv4RootNode)
	if err != nil {
		return err
	}

	for _, network := range ipv4AliasNetworks {
		err := t.insertStringNetwork(network, recordTypeAlias, nil, ipv4RootNode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tree) insertReservedNetworks() error {
	// the reserved networks are in reserved.go
	networks := reservedNetworksIPv4
	if t.ipVersion == 6 {
		networks = append(networks, reservedNetworksIPv6...)
	}

	for _, network := range networks {
		err := t.insertStringNetwork(network, recordTypeReserved, nil, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get the value for the given IP address from the tree. If the nil interface
// is returned, that means the tree does not have a value for the IP.
func (t *Tree) Get(ip net.IP) (*net.IPNet, mmdbtype.DataType) {
	lookupIP := ip

	if t.treeDepth == 128 {
		// We use To4() here as Go will parse an IPv4 address to a 16 byte
		// IPv6-mapped IPv4 address, e.g.:
		//
		// len(net.ParseIP("1.1.1.1")) == 16
		//
		// The parsed address above is equal to ::ffff:1.1.1.1. However,
		// the MaxMind DB format has the record for 1.1.1.1 at ::1.1.1.1.
		if ipv4 := ip.To4(); ipv4 != nil {
			lookupIP = ipV4ToV6(ipv4)
		}
	}

	prefixLen, r := t.root.get(lookupIP, 0)

	// This is so that if you look up an IPv4 address in a database that has
	// an IPv4 subtree, you will get back an IPv4 network. This matches what
	// github.com/oschwald/maxminddb-golang does.
	if prefixLen >= 96 && len(ip) == 4 {
		prefixLen -= 96
	}

	mask := net.CIDRMask(prefixLen, t.treeDepth)

	var value mmdbtype.DataType
	if r.recordType == recordTypeData {
		value = r.value.data
	}

	return &net.IPNet{
		IP:   ip.Mask(mask),
		Mask: mask,
	}, value
}

// finalize prepares the tree for writing. It is not threadsafe.
func (t *Tree) finalize() {
	_, t.nodeCount = t.root.finalize(0)
}

// WriteTo writes the tree to the provided Writer.
func (t *Tree) WriteTo(w io.Writer) (int64, error) {
	if t.nodeCount == 0 {
		t.finalize()
	}

	buf := bufio.NewWriter(w)
	//nolint:errcheck // We check the error on flush the only place that matters.
	defer buf.Flush()

	// We create this here so that we don't have to allocate millions of these. This
	// may no longer make sense now that we are using a bufio.Writer anyway, which has
	// WriteByte, but we should probably do some testing.
	recordBuf := make([]byte, 2*t.recordSize/8)

	usePointers := true
	dataWriter := newDataWriter(t.dataMap, usePointers)

	nodeCount, numBytes, err := t.writeNode(buf, t.root, dataWriter, recordBuf)
	if err != nil {
		return numBytes, err
	}
	if nodeCount != t.nodeCount {
		// This should only happen if there is a programming bug
		// in this library.
		return numBytes, fmt.Errorf(
			"number of nodes written (%d) doesn't match number expected (%d)",
			nodeCount,
			t.nodeCount,
		)
	}

	nb, err := buf.Write(dataSectionSeparator)
	numBytes += int64(nb)
	if err != nil {
		return numBytes, fmt.Errorf("writing data section separator: %w", err)
	}

	nb64, err := dataWriter.WriteTo(buf)
	numBytes += nb64
	if err != nil {
		return numBytes, err
	}

	nb, err = buf.Write(metadataStartMarker)
	numBytes += int64(nb)
	if err != nil {
		return numBytes, fmt.Errorf("writing metadata start marker: %w", err)
	}

	metadataWriter := newDataWriter(dataWriter.dataMap, !t.disableMetadataPointers)
	_, err = t.writeMetadata(metadataWriter)
	if err != nil {
		return numBytes, fmt.Errorf("writing metadata: %w", err)
	}

	nb64, err = metadataWriter.WriteTo(buf)
	numBytes += nb64
	if err != nil {
		return numBytes, fmt.Errorf("writing metadata to buffer: %w", err)
	}

	err = buf.Flush()
	if err != nil {
		return numBytes, fmt.Errorf("flushing buffer to writer: %w", err)
	}

	return numBytes, err
}

func (t *Tree) writeNode(
	w io.Writer,
	n *node,
	dataWriter *dataWriter,
	recordBuf []byte,
) (int, int64, error) {
	err := t.copyNode(recordBuf, n, dataWriter)
	if err != nil {
		return 0, 0, err
	}

	numBytes := int64(0)
	nb, err := w.Write(recordBuf)
	numBytes += int64(nb)
	nodesWritten := 1
	if err != nil {
		return nodesWritten, numBytes, fmt.Errorf("writing node: %w", err)
	}

	for i := 0; i < 2; i++ {
		child := n.children[i]
		if child.recordType != recordTypeNode && child.recordType != recordTypeFixedNode {
			continue
		}
		addedNodes, addedBytes, err := t.writeNode(
			w,
			n.children[i].node,
			dataWriter,
			recordBuf,
		)
		nodesWritten += addedNodes
		numBytes += addedBytes
		if err != nil {
			return nodesWritten, numBytes, err
		}
	}

	return nodesWritten, numBytes, nil
}

func (t *Tree) recordValue(
	r record,
	dataWriter *dataWriter,
) (int, error) {
	switch r.recordType {
	case recordTypeData:
		offset, err := dataWriter.maybeWrite(r.value)
		return t.nodeCount + len(dataSectionSeparator) + offset, err
	case recordTypeEmpty, recordTypeReserved:
		return t.nodeCount, nil
	default:
		return r.node.nodeNum, nil
	}
}

func (t *Tree) copyNode(buf []byte, n *node, dataWriter *dataWriter) error {
	left, err := t.recordValue(n.children[0], dataWriter)
	if err != nil {
		return err
	}
	right, err := t.recordValue(n.children[1], dataWriter)
	if err != nil {
		return err
	}

	maxRecord := 1 << t.recordSize
	if left >= maxRecord || right >= maxRecord {
		return fmt.Errorf(
			"exceeded record capacity by attempting to write (%d, %d) to node with %d bit record size; "+
				"try increasing RecordSize or reducing the size of the database",
			left,
			right,
			t.recordSize,
		)
	}

	switch t.recordSize {
	case 24:
		buf[0] = byte((left >> 16) & 0xFF)
		buf[1] = byte((left >> 8) & 0xFF)
		buf[2] = byte(left & 0xFF)
		buf[3] = byte((right >> 16) & 0xFF)
		buf[4] = byte((right >> 8) & 0xFF)
		buf[5] = byte(right & 0xFF)
	case 28:
		buf[0] = byte((left >> 16) & 0xFF)
		buf[1] = byte((left >> 8) & 0xFF)
		buf[2] = byte(left & 0xFF)
		buf[3] = byte((((left >> 24) & 0x0F) << 4) | (right >> 24 & 0x0F))
		buf[4] = byte((right >> 16) & 0xFF)
		buf[5] = byte((right >> 8) & 0xFF)
		buf[6] = byte(right & 0xFF)
	case 32:
		buf[0] = byte((left >> 24) & 0xFF)
		buf[1] = byte((left >> 16) & 0xFF)
		buf[2] = byte((left >> 8) & 0xFF)
		buf[3] = byte(left & 0xFF)
		buf[4] = byte((right >> 24) & 0xFF)
		buf[5] = byte((right >> 16) & 0xFF)
		buf[6] = byte((right >> 8) & 0xFF)
		buf[7] = byte(right & 0xFF)
	default:
		return fmt.Errorf("unsupported record size of %d", t.recordSize)
	}
	return nil
}

var v4Prefix = net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func ipV4ToV6(ip net.IP) net.IP {
	return append(v4Prefix, ip...)
}

func (t *Tree) writeMetadata(dw *dataWriter) (int64, error) {
	description := mmdbtype.Map{}
	for k, v := range t.description {
		description[mmdbtype.String(k)] = mmdbtype.String(v)
	}

	languages := mmdbtype.Slice{}
	for _, v := range t.languages {
		languages = append(languages, mmdbtype.String(v))
	}
	metadata := mmdbtype.Map{
		"binary_format_major_version": mmdbtype.Uint16(2),
		"binary_format_minor_version": mmdbtype.Uint16(0),
		"build_epoch":                 mmdbtype.Uint64(t.buildEpoch),
		"database_type":               mmdbtype.String(t.databaseType),
		"description":                 description,
		"ip_version":                  mmdbtype.Uint16(t.ipVersion),
		"languages":                   languages,
		"node_count":                  mmdbtype.Uint32(t.nodeCount),
		"record_size":                 mmdbtype.Uint16(t.recordSize),
	}
	return metadata.WriteTo(dw)
}
