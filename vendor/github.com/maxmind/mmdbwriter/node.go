package mmdbwriter

import (
	"fmt"
	"net"

	"github.com/maxmind/mmdbwriter/mmdbtype"
)

type recordType byte

const (
	recordTypeEmpty recordType = iota
	recordTypeData
	recordTypeNode
	recordTypeAlias
	recordTypeFixedNode
	recordTypeReserved
)

type record struct {
	node       *node
	value      *dataMapValue
	recordType recordType
}

// each node contains two records.
type node struct {
	children [2]record
	nodeNum  int
}

type insertRecord struct {
	inserter func(value mmdbtype.DataType) (mmdbtype.DataType, error)

	dataMap      *dataMap
	insertedNode *node

	ip        net.IP
	prefixLen int

	recordType recordType
}

func (n *node) insert(iRec insertRecord, currentDepth int) error {
	newDepth := currentDepth + 1
	// Check if we are inside the network already
	if newDepth > iRec.prefixLen {
		// Data already exists for the network so insert into all the children.
		// We will prune duplicate nodes when we finalize.
		err := n.children[0].insert(iRec, newDepth)
		if err != nil {
			return err
		}
		return n.children[1].insert(iRec, newDepth)
	}

	// We haven't reached the network yet.
	pos := bitAt(iRec.ip, currentDepth)
	r := &n.children[pos]
	return r.insert(iRec, newDepth)
}

func (r *record) insert(
	iRec insertRecord,
	newDepth int,
) error {
	switch r.recordType {
	case recordTypeNode, recordTypeFixedNode:
	case recordTypeEmpty, recordTypeData:
		if newDepth >= iRec.prefixLen {
			r.node = iRec.insertedNode
			r.recordType = iRec.recordType
			if iRec.recordType == recordTypeData {
				var data mmdbtype.DataType
				if r.value != nil {
					data = r.value.data

					// Potentially we could avoid this if the
					// new value is the same, but it would likely
					// not save us much and the code would be a
					// bit more complicated.
					iRec.dataMap.remove(r.value)
				}
				value, err := iRec.inserter(data)
				if err != nil {
					return err
				}
				if value == nil {
					r.recordType = recordTypeEmpty
					r.value = nil
				} else {
					value, err := iRec.dataMap.store(value)
					if err != nil {
						return err
					}
					r.value = value
				}
			} else {
				r.value = nil
			}
			return nil
		}

		// We are splitting this record so we create two duplicate child
		// records.
		r.node = &node{children: [2]record{*r, *r}}
		r.value = nil
		r.recordType = recordTypeNode
	case recordTypeReserved:
		if iRec.prefixLen >= newDepth {
			return fmt.Errorf(
				"attempt to insert %s/%d, which is in a reserved network",
				iRec.ip,
				iRec.prefixLen,
			)
		}
		// If we are inserting a network that contains a reserved network,
		// we silently remove the reserved network.
		return nil
	case recordTypeAlias:
		if iRec.prefixLen < newDepth {
			// Do nothing. We are inserting a network that contains an aliased
			// network. We silently ignore.
			return nil
		}
		// attempting to insert _into_ an aliased network
		return fmt.Errorf(
			"attempt to insert %s/%d, which is in an aliased network",
			iRec.ip,
			iRec.prefixLen,
		)
	default:
		return fmt.Errorf("inserting into record type %d is not implemented", r.recordType)
	}

	return r.node.insert(iRec, newDepth)
}

func (n *node) get(
	ip net.IP,
	depth int,
) (int, record) {
	r := n.children[bitAt(ip, depth)]

	depth++

	switch r.recordType {
	case recordTypeNode, recordTypeAlias, recordTypeFixedNode:
		return r.node.get(ip, depth)
	default:
		return depth, r
	}
}

// finalize prunes unnecessary nodes (e.g., where the two records are the same) and
// sets the node number for the node. It returns a record pointer that is nil if
// the node is not mergeable or the value of the merged record if it can be merged.
// The second return value is the current node count, including the subtree.
func (n *node) finalize(currentNum int) (*record, int) {
	n.nodeNum = currentNum
	currentNum++

	for i := 0; i < 2; i++ {
		switch n.children[i].recordType {
		case recordTypeFixedNode:
			// We don't consider merging for fixed nodes
			_, currentNum = n.children[i].node.finalize(currentNum)
		case recordTypeNode:
			record, newCurrentNum := n.children[i].node.finalize(currentNum)
			if record == nil {
				// nothing to merge. Use current number from child.
				currentNum = newCurrentNum
			} else {
				n.children[i] = *record
			}
		default:
		}
	}

	if n.children[0].recordType == n.children[1].recordType &&
		(n.children[0].recordType == recordTypeEmpty ||
			(n.children[0].recordType == recordTypeData &&
				n.children[0].value.key == n.children[1].value.key)) {
		return &record{
			recordType: n.children[0].recordType,
			value:      n.children[0].value,
		}, currentNum
	}

	return nil, currentNum
}

func bitAt(ip net.IP, depth int) byte {
	return (ip[depth/8] >> (7 - (depth % 8))) & 1
}
