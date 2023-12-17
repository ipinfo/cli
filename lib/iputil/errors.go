package iputil

import (
	"errors"
)

var (
	// ErrNotASN is returned when the input is not in proper ASN form.
	ErrNotASN = errors.New("not asn")

	// ErrNotIP is returned when the input is not in proper IP form.
	ErrNotIP = errors.New("not ip")

	// ErrNotCIDR is returned when the input is not in proper CIDR form.
	ErrNotCIDR = errors.New("not cidr")

	// ErrNotIPRange is returned when the input is not in proper IP range form.
	ErrNotIPRange = errors.New("not ip range")

	// ErrNotIP6Range is returned when the input is not in proper IPv6 range
	// form.
	ErrNotIP6Range = errors.New("not ipv6 range")

	// ErrNotFile is returned when the input is not an actual file.
	ErrNotFile = errors.New("not file")

	// ErrMissingCIDRsOrIPRange is returned when a CIDR or IP range is missing.
	ErrMissingCIDRsOrIPRange = errors.New("missing CIDRs or IP range")

	// ErrCannotMixCIDRAndIPs is returned when CIDRs and IPs are mixed.
	ErrCannotMixCIDRAndIPs = errors.New("cannot mix CIDRs and IPs")

	// ErrIPRangeRequiresTwoIPs is returned when the IP range is incomplete.
	ErrIPRangeRequiresTwoIPs = errors.New("IP range requires 2 IPs")

	// ErrInvalidInput is returned as a generic error for bad input.
	ErrInvalidInput = errors.New("invalid input")
)
