package ipinfo

import (
	"net/netip"
)

func isBogon(ip netip.Addr) bool {
	for _, network := range bogonNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

var bogonNetworks = []netip.Prefix{
	netip.MustParsePrefix("0.0.0.0/8"),
	netip.MustParsePrefix("10.0.0.0/8"),
	netip.MustParsePrefix("100.64.0.0/10"),
	netip.MustParsePrefix("127.0.0.0/8"),
	netip.MustParsePrefix("169.254.0.0/16"),
	netip.MustParsePrefix("172.16.0.0/12"),
	netip.MustParsePrefix("192.0.0.0/24"),
	netip.MustParsePrefix("192.0.2.0/24"),
	netip.MustParsePrefix("192.168.0.0/16"),
	netip.MustParsePrefix("198.18.0.0/15"),
	netip.MustParsePrefix("198.51.100.0/24"),
	netip.MustParsePrefix("203.0.113.0/24"),
	netip.MustParsePrefix("224.0.0.0/4"),
	netip.MustParsePrefix("240.0.0.0/4"),
	netip.MustParsePrefix("255.255.255.255/32"),
	netip.MustParsePrefix("::/128"),
	netip.MustParsePrefix("::1/128"),
	netip.MustParsePrefix("::ffff:0:0/96"),
	netip.MustParsePrefix("::/96"),
	netip.MustParsePrefix("100::/64"),
	netip.MustParsePrefix("2001:10::/28"),
	netip.MustParsePrefix("2001:db8::/32"),
	netip.MustParsePrefix("fc00::/7"),
	netip.MustParsePrefix("fe80::/10"),
	netip.MustParsePrefix("fec0::/10"),
	netip.MustParsePrefix("ff00::/8"),
	netip.MustParsePrefix("2002::/24"),
	netip.MustParsePrefix("2002:a00::/24"),
	netip.MustParsePrefix("2002:7f00::/24"),
	netip.MustParsePrefix("2002:a9fe::/32"),
	netip.MustParsePrefix("2002:ac10::/28"),
	netip.MustParsePrefix("2002:c000::/40"),
	netip.MustParsePrefix("2002:c000:200::/40"),
	netip.MustParsePrefix("2002:c0a8::/32"),
	netip.MustParsePrefix("2002:c612::/31"),
	netip.MustParsePrefix("2002:c633:6400::/40"),
	netip.MustParsePrefix("2002:cb00:7100::/40"),
	netip.MustParsePrefix("2002:e000::/20"),
	netip.MustParsePrefix("2002:f000::/20"),
	netip.MustParsePrefix("2002:ffff:ffff::/48"),
	netip.MustParsePrefix("2001::/40"),
	netip.MustParsePrefix("2001:0:a00::/40"),
	netip.MustParsePrefix("2001:0:7f00::/40"),
	netip.MustParsePrefix("2001:0:a9fe::/48"),
	netip.MustParsePrefix("2001:0:ac10::/44"),
	netip.MustParsePrefix("2001:0:c000::/56"),
	netip.MustParsePrefix("2001:0:c000:200::/56"),
	netip.MustParsePrefix("2001:0:c0a8::/48"),
	netip.MustParsePrefix("2001:0:c612::/47"),
	netip.MustParsePrefix("2001:0:c633:6400::/56"),
	netip.MustParsePrefix("2001:0:cb00:7100::/56"),
	netip.MustParsePrefix("2001:0:e000::/36"),
	netip.MustParsePrefix("2001:0:f000::/36"),
	netip.MustParsePrefix("2001:0:ffff:ffff::/64"),
}
