package ipUtils

import (
	"regexp"
)

func init() {
	BogonIP4List = GetBogonRange4()
	BogonIP6List = GetBogonRange6()

	// IP
	IpV4Regex = regexp.MustCompilePOSIX(IPv4RegexPattern)
	IpV6Regex = regexp.MustCompilePOSIX(IPv6RegexPattern)
	IpRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern)

	// IP and CIDR
	V4IpCidrRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv4CIDRRegexPattern)
	V6IpCidrRegex = regexp.MustCompilePOSIX(IPv6RegexPattern + "|" + IPv6CIDRRegexPattern)
	IpCidrRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern + "|" + IPv4CIDRRegexPattern + "|" + IPv6CIDRRegexPattern)

	// IP and Range
	V4IpRangeRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv4RangeRegexPattern)
	V6IpRangeRegex = regexp.MustCompilePOSIX(IPv6RegexPattern + "|" + IPv6RangeRegexPattern)
	IpRangeRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern + "|" + IPv4RangeRegexPattern + "|" + IPv6RangeRegexPattern)

	// IP, CIDR and Range
	V4IpSubnetRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv4CIDRRegexPattern + "|" + IPv4RangeRegexPattern)
	V6IpSubnetRegex = regexp.MustCompilePOSIX(IPv6RegexPattern + "|" + IPv6RangeRegexPattern + "|" + IPv6CIDRRegexPattern)
	IpSubnetRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern + "|" + IPv4CIDRRegexPattern + "|" + IPv4RangeRegexPattern + "|" + IPv6RangeRegexPattern + "|" + IPv6CIDRRegexPattern)

	// CIDR
	V4CidrRegex = regexp.MustCompilePOSIX(IPv4CIDRRegexPattern)
	V6CidrRegex = regexp.MustCompilePOSIX(IPv6CIDRRegexPattern)
	CidrRegex = regexp.MustCompilePOSIX(IPv4CIDRRegexPattern + "|" + IPv6CIDRRegexPattern)

	// Range
	V4RangeRegex = regexp.MustCompilePOSIX(IPv4RangeRegexPattern)
	V6RangeRegex = regexp.MustCompilePOSIX(IPv6RangeRegexPattern)
	RangeRegex = regexp.MustCompilePOSIX(IPv4RangeRegexPattern + "|" + IPv6RangeRegexPattern)

	// CIDR and Range
	V4SubnetRegex = regexp.MustCompilePOSIX(IPv4RangeRegexPattern + "|" + IPv4CIDRRegexPattern)
	V6SubnetRegex = regexp.MustCompilePOSIX(IPv6RangeRegexPattern + "|" + IPv6CIDRRegexPattern)
	SubnetRegex = regexp.MustCompilePOSIX(IPv4RangeRegexPattern + "|" + IPv4CIDRRegexPattern + "|" + IPv6RangeRegexPattern + "|" + IPv6CIDRRegexPattern)
}
