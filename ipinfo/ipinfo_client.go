package main

import (
	"fmt"
	"runtime"

	"github.com/ipinfo/go/v2/ipinfo"
)

var ii *ipinfo.Client

func prepareIpinfoClient(tok string) *ipinfo.Client {
	var _ii *ipinfo.Client

	// get token from persistent store.
	if tok == "" {
		tok = restoreToken()
	}

	var cache *ipinfo.Cache
	if gConfig.Cache && !fNoCache {
		boltdbCache, err := NewBoltdbCache()
		if err != nil {
			fmt.Printf("warn: cache will not be used: %v", err)
		} else {
			cache = ipinfo.NewCache(boltdbCache)
		}
	}

	// init client.
	_ii = ipinfo.NewClient(nil, cache, tok)
	_ii.UserAgent = fmt.Sprintf(
		"IPinfoCli/%s (os/%s - arch/%s)",
		version, runtime.GOOS, runtime.GOARCH,
	)
	return _ii
}
