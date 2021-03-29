package lib

import (
	"errors"
)

var (
	ErrNotASN                = errors.New("not asn")
	ErrNotIP                 = errors.New("not ip")
	ErrNotCIDR               = errors.New("not cidr")
	ErrNotIPRange            = errors.New("not ip range")
	ErrNotFile               = errors.New("not file")
	ErrMissingCIDRsOrIPRange = errors.New("missing CIDRs or IP range")
	ErrCannotMixCIDRAndIPs   = errors.New("cannot mix CIDRs and IPs")
	ErrIPRangeRequiresTwoIPs = errors.New("IP range requires 2 IPs")
	ErrInvalidInput          = errors.New("invalid input")
)
