package main

import (
	"errors"
)

var (
	errNotASN                = errors.New("not asn")
	errNotIP                 = errors.New("not ip")
	errMissingCIDRsOrIPRange = errors.New("missing CIDRs or IP range")
	errCannotMixCIDRAndIPs   = errors.New("cannot mix CIDRs and IPs")
	errIPRangeRequiresTwoIPs = errors.New("IP range requires 2 IPs")
)
