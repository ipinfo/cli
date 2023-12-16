package ipUtils

// these lists are initialized on startup inside this pkg's `init`.
var BogonIP4List []IPRange
var BogonIP6List []IP6Range

// list of bogon IPv4 IPs.
var BogonRange4Str []string = []string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"192.168.0.0/16",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"224.0.0.0/4",
	"240.0.0.0/4",
	"255.255.255.255/32",
}

// list of bogon IPv6 IPs.
var BogonRange6Str []string = []string{
	"::/128",
	"::1/128",
	"::ffff:0:0/96",
	"::/96",
	"100::/64",
	"2001:10::/28",
	"2001:db8::/32",
	"fc00::/7",
	"fe80::/10",
	"fec0::/10",
	"ff00::/8",
	// 6to4 bogon ranges
	"2002::/24",
	"2002:a00::/24",
	"2002:7f00::/24",
	"2002:a9fe::/32",
	"2002:ac10::/28",
	"2002:c000::/40",
	"2002:c000:200::/40",
	"2002:c0a8::/32",
	"2002:c612::/31",
	"2002:c633:6400::/40",
	"2002:cb00:7100::/40",
	"2002:e000::/20",
	"2002:f000::/20",
	"2002:ffff:ffff::/48",
	// teredo
	"2001::/40",
	"2001:0:a00::/40",
	"2001:0:7f00::/40",
	"2001:0:a9fe::/48",
	"2001:0:ac10::/44",
	"2001:0:c000::/56",
	"2001:0:c000:200::/56",
	"2001:0:c0a8::/48",
	"2001:0:c612::/47",
	"2001:0:c633:6400::/56",
	"2001:0:cb00:7100::/56",
	"2001:0:e000::/36",
	"2001:0:f000::/36",
	"2001:0:ffff:ffff::/64",
}

// GetBogonRange4 returns list of IPRange of all IPv4 bogon IPs.
func GetBogonRange4() []IPRange {
	bogonRanges4 := make([]IPRange, len(BogonRange4Str))
	for i, bogonRangeStr := range BogonRange4Str {
		r, err := IPRangeFromCIDR(bogonRangeStr)
		if err != nil {
			panic(err)
		}

		bogonRanges4[i] = r
	}
	return bogonRanges4
}

// IsBogonIP4 returns true if IPv4 is a BogonIP.
func IsBogonIP4(ip uint32) bool {
	for _, bogonIP := range BogonIP4List {
		if ip >= uint32(bogonIP.Start) && ip <= uint32(bogonIP.End) {
			return true
		}
	}
	return false
}

// GetBogonRange6 returns list of IPRange of all IPv6 bogon IPs.
func GetBogonRange6() []IP6Range {
	bogonRanges6 := make([]IP6Range, len(BogonRange6Str))
	for i, bogonRangeStr := range BogonRange6Str {
		r, err := IP6RangeFromCIDR(bogonRangeStr)
		if err != nil {
			panic(err)
		}

		bogonRanges6[i] = r
	}
	return bogonRanges6
}

// IsBogonIP6 returns true if IPv6 is a BogonIP.
func IsBogonIP6(ip6 U128) bool {
	for _, bogonIP := range BogonIP6List {
		if (ip6.Cmp(bogonIP.Start.N) >= 0) && (ip6.Cmp(bogonIP.End.N) <= 0) {
			return true
		}
	}
	return false
}
