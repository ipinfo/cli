package mmdbwriter

// These were taken from the Perl writer.
//
// https://www.iana.org/assignments/iana-ipv4-special-registry/iana-ipv4-special-registry.xhtml
var reservedNetworksIPv4 = []string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	// This is an odd case. 192.0.0.0/24 is reserved, but there is a note that
	// says "Not useable unless by virtue of a more specific reservation". As
	// such, since 192.0.0.0/29 was more recently reserved, it's possible the
	// intention is that the rest is not reserved any longer. I'm not too clear
	// on this, but I believe that is the rationale, so I choose to leave it.
	"192.0.0.0/29",
	// TODO(wstorey@maxmind.com): 192.168.0.8/32
	// TODO(wstorey@maxmind.com): 192.168.0.9/32
	// TODO(wstorey@maxmind.com): 192.168.0.10/32
	// TODO(wstorey@maxmind.com): 192.168.0.170/32
	// TODO(wstorey@maxmind.com): 192.168.0.171/32
	"192.0.2.0/24",
	// 192.31.196.0/24 is routable I believe
	// TODO(wstorey@maxmnd.com): 192.52.193.0/24
	// TODO(wstorey@maxmind.com): Looks like 192.88.99.0/24 may no longer be
	// reserved?
	"192.88.99.0/24",
	"192.168.0.0/16",
	// 192.175.48.0/24 is routable I believe
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	// The above IANA page doesn't list 224.0.0.0/4, but at least some parts
	// are listed in https://tools.ietf.org/html/rfc5771
	"224.0.0.0/4",
	"240.0.0.0/4",
	// 255.255.255.255/32 gets brought in by 240.0.0.0/4.
}

// https://www.iana.org/assignments/iana-ipv6-special-registry/iana-ipv6-special-registry.xhtml
var reservedNetworksIPv6 = []string{
	// ::/128 and ::1/128 are reserved under IPv6 but these are already
	// covered under 0.0.0.0/8.
	//
	// ::ffff:0:0/96 - IPv4 mapped addresses. We treat it specially with the
	// `alias_ipv6_to_ipv4' option.
	//
	// 64:ff9b::/96 - well known prefix mapping, covered by alias_ipv6_to_ipv4
	//
	// TODO(wstorey@maxmind.com): 64:ff9b:1::/48 should be in
	// alias_ipv6_to_ipv4?

	"100::/64",

	// 2001::/23 is reserved. We include all of it here other than 2001::/32
	// as it is Teredo which is globally routable.
	"2001:1::/32",
	"2001:2::/31",
	"2001:4::/30",
	"2001:8::/29",
	"2001:10::/28",
	"2001:20::/27",
	"2001:40::/26",
	"2001:80::/25",
	"2001:100::/24",

	"2001:db8::/32",
	// 2002::/16 - 6to4, part of alias_ipv6_to_ipv4
	// 2620:4f:8000::/48 is routable I believe
	"fc00::/7",
	"fe80::/10",
	// Multicast
	"ff00::/8",
}
