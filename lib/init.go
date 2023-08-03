package lib

import "regexp"

func init() {
	bogonIP4List = GetBogonRange4()
	bogonIP6List = GetBogonRange6()
	ipV4Regex = regexp.MustCompile(IPv4RegexPattern)
	ipV6Regex = regexp.MustCompile(IPv6RegexPattern)
	ipRegex = regexp.MustCompile(IPv4RegexPattern + "|" + IPv6RegexPattern)
}
