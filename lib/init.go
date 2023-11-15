package lib

import "regexp"

func init() {
	bogonIP4List = GetBogonRange4()
	bogonIP6List = GetBogonRange6()

	// IP
	ipV4Regex = regexp.MustCompilePOSIX(IPv4RegexPattern)
	ipV6Regex = regexp.MustCompilePOSIX(IPv6RegexPattern)
	ipRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern)

	// IP and CIDR
	v4IpCidrRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv4CIDRRegexPattern)
	v6IpCidrRegex = regexp.MustCompilePOSIX(IPv6RegexPattern + "|" + IPv6CIDRRegexPattern)
	ipCidrRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern + "|" + IPv4CIDRRegexPattern + "|" + IPv6CIDRRegexPattern)

	// IP and Range
	v4IpRangeRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv4RangeRegexPattern)
	v6IpRangeRegex = regexp.MustCompilePOSIX(IPv6RegexPattern + "|" + IPv6RangeRegexPattern)
	ipRangeRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern + "|" + IPv4RangeRegexPattern + "|" + IPv6RangeRegexPattern)

	// IP, CIDR and Range
	v4IpSubnetRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv4CIDRRegexPattern + "|" + IPv4RangeRegexPattern)
	v6IpSubnetRegex = regexp.MustCompilePOSIX(IPv6RegexPattern + "|" + IPv6RangeRegexPattern + "|" + IPv6CIDRRegexPattern)
	ipSubnetRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern + "|" + IPv4CIDRRegexPattern + "|" + IPv4RangeRegexPattern + "|" + IPv6RangeRegexPattern + "|" + IPv6CIDRRegexPattern)

	// CIDR
	v4CidrRegex = regexp.MustCompilePOSIX(IPv4CIDRRegexPattern)
	v6CidrRegex = regexp.MustCompilePOSIX(IPv6CIDRRegexPattern)
	cidrRegex = regexp.MustCompilePOSIX(IPv4CIDRRegexPattern + "|" + IPv6CIDRRegexPattern)

	// Range
	v4RangeRegex = regexp.MustCompilePOSIX(IPv4RangeRegexPattern)
	v6RangeRegex = regexp.MustCompilePOSIX(IPv6RangeRegexPattern)
	rangeRegex = regexp.MustCompilePOSIX(IPv4RangeRegexPattern + "|" + IPv6RangeRegexPattern)

	// CIDR and Range
	v4SubnetRegex = regexp.MustCompilePOSIX(IPv4RangeRegexPattern + "|" + IPv4CIDRRegexPattern)
	v6SubnetRegex = regexp.MustCompilePOSIX(IPv6RangeRegexPattern + "|" + IPv6CIDRRegexPattern)
	subnetRegex = regexp.MustCompilePOSIX(IPv4RangeRegexPattern + "|" + IPv4CIDRRegexPattern + "|" + IPv6RangeRegexPattern + "|" + IPv6CIDRRegexPattern)
}
