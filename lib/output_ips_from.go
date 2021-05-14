package lib

// OutputIPsFrom outputs a list of IPs from inputs which are interpreted to
// contain IP ranges and IP CIDRs in them, all depending upon which flags are
// set.
func OutputIPsFrom(
	inputs []string,
	iprange bool,
	cidr bool,
) error {
	// prevent edge cases with all flags turned off.
	if !iprange && !cidr {
		return nil
	}

	for _, input := range inputs {
		if iprange {
			if err := OutputIPsFromRangeStr(input); err == nil {
				continue
			}
		}

		if cidr && IsCIDR(input) {
			if err := OutputIPsFromCIDR(input); err == nil {
				continue
			}
		}

		return ErrInvalidInput
	}

	return nil
}
