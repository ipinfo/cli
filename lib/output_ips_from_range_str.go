package lib

// OutputIPsFromRangeStr outputs all IPs in an IP range.
// The string must be of any of these forms:
//	<ip_range_start>-<ip_range_end>
//	<ip_range_start>,<ip_range_end>
func OutputIPsFromRangeStr(r string) error {
	rStart, rEnd, err := IPRangeStrPartsFromRangeStr(r)
	if err != nil {
		return err
	}

	return OutputIPsFromRange(rStart, rEnd)
}
