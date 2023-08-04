package lib

import "regexp"

func init() {
	bogonIP4List = GetBogonRange4()
	bogonIP6List = GetBogonRange6()
	ipV4Regex = regexp.MustCompilePOSIX(IPv4RegexPattern)
	ipV6Regex = regexp.MustCompilePOSIX(IPv6RegexPattern)
	ipRegex = regexp.MustCompilePOSIX(IPv4RegexPattern + "|" + IPv6RegexPattern)
}
