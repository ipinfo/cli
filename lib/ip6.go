package lib

// IP6u128 is a 128-bit number representation of a IPv6 address in big endian
// byte order.
//
// The number is internally represented as 2 64-bit numbers.
type IP6u128 struct {
	Hi uint64
	Lo uint64
}

// Cmp compares `ip1` and `ip2` and returns:
//
// - -1 if `ip1<ip2`
// -  0 if `ip1==ip2`
// -  1 if `ip1>ip2`
func (ip1 IP6u128) Cmp(ip2 IP6u128) int {
	if ip1 == ip2 {
		return 0
	} else if ip1.Hi < ip2.Hi || (ip1.Hi == ip2.Hi && ip1.Lo < ip2.Lo) {
		return -1
	} else {
		return 1
	}
}

// Eq returns whether `ip1 == ip2`.
func (ip1 IP6u128) Eq(ip2 IP6u128) bool {
	return ip1.Cmp(ip2) == 0
}

// Gt returns whether `ip1 > ip2`.
func (ip1 IP6u128) Gt(ip2 IP6u128) bool {
	return ip1.Cmp(ip2) > 0
}

// Gte returns whether `ip1 >= ip2`.
func (ip1 IP6u128) Gte(ip2 IP6u128) bool {
	return ip1.Cmp(ip2) >= 0
}

// Lt returns whether `ip1 < ip2`.
func (ip1 IP6u128) Lt(ip2 IP6u128) bool {
	return ip1.Cmp(ip2) < 0
}

// Lte returns whether `ip1 <= ip2`.
func (ip1 IP6u128) Lte(ip2 IP6u128) bool {
	return ip1.Cmp(ip2) <= 0
}
