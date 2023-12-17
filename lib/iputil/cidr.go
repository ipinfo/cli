package iputil

import (
	"bytes"
	"encoding/binary"
	"math"
	"net"
	"sort"
)

// CIDR represens a Classless Inter-Domain Routing structure.
type CIDR struct {
	IP      net.IP
	Network *net.IPNet
}

// NewCidr creates a NewCidr CIDR structure.
func NewCidr(s string) *CIDR {
	ip, ipnet, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return &CIDR{
		IP:      ip,
		Network: ipnet,
	}
}

func (c *CIDR) String() string {
	return c.Network.String()
}

// MaskLen returns a network mask length.
func (c *CIDR) MaskLen() uint32 {
	i, _ := c.Network.Mask.Size()
	return uint32(i)
}

// PrefixUint32 returns a prefix.
func (c *CIDR) PrefixUint32() uint32 {
	return binary.BigEndian.Uint32(c.IP.To4())
}

// Size returns a size of a CIDR range.
func (c *CIDR) Size() int {
	ones, bits := c.Network.Mask.Size()
	return int(math.Pow(2, float64(bits-ones)))
}

// NewCidrList returns a slice of sorted CIDR structures.
func NewCidrList(s []string) []*CIDR {
	out := make([]*CIDR, 0)
	for _, c := range s {
		out = append(out, NewCidr(c))
	}
	sort.Sort(cidrSort(out))
	return out
}

type cidrSort []*CIDR

func (s cidrSort) Len() int      { return len(s) }
func (s cidrSort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s cidrSort) Less(i, j int) bool {
	cmp := bytes.Compare(s[i].IP, s[j].IP)
	return cmp < 0 || (cmp == 0 && s[i].MaskLen() < s[j].MaskLen())
}
