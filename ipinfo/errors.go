package main

import (
	"errors"
)

var (
	errNotASN                = errors.New("not asn")
	errNotIP                 = errors.New("not ip")
	errNotCIDR               = errors.New("not cidr")
	errNotIPRange            = errors.New("not ip range")
	errNotFile               = errors.New("not file")
	errMissingCIDRsOrIPRange = errors.New("missing CIDRs or IP range")
	errCannotMixCIDRAndIPs   = errors.New("cannot mix CIDRs and IPs")
	errIPRangeRequiresTwoIPs = errors.New("IP range requires 2 IPs")
	errInvalidInput          = errors.New("invalid input")
)
