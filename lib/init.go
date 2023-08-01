package lib

func init() {
	bogonIP4List = GetBogonRange4()
	bogonIP6List = GetBogonRange6()
	ipV4RgxPattern = GetIpV4RgxPattern()
	ipV6RgxPattern = GetIpV6RgxPattern()
}
