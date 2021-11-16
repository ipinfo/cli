package lib

var bogonIP4List []IPRange
var bogonIP6List []IP6Range

func init() {
	bogonIP4List = GetBogonRange4()
	bogonIP6List = GetBogonRange6()

}
