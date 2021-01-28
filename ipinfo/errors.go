package main

import (
	"errors"
)

var (
	errNotASN = errors.New("not asn")
	errNotIP  = errors.New("not ip")
)
