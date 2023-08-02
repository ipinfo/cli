package ipinfo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	batchMaxSize                        = 1000
	batchReqTimeoutDefault              = 5
	batchDefaultConcurrentRequestsLimit = 8
)

// Internal batch type used by common batch functionality to temporarily store
// the URL-to-result mapping in a half-decoded state (specifically the value
// not being decoded yet). This allows us to decode the value to a proper
// concrete type like `Core` or `ASNDetails` after analyzing the key to
// determine which one it should be.
type batch map[string]json.RawMessage

// Batch is a mapped result of any valid API endpoint (e.g. `<ip>`,
// `<ip>/<field>`, `<asn>`, etc) to its corresponding data.
//
// The corresponding value will be either `*Core`, `*ASNDetails` or a generic
// map for unknown value results.
type Batch map[string]interface{}

// BatchCore is a mapped result of IPs to their corresponding `Core` data.
type BatchCore map[string]*Core

// BatchASNDetails is a mapped result of ASNs to their corresponding
// `ASNDetails` data.
type BatchASNDetails map[string]*ASNDetails

// BatchReqOpts are options input into batch request functions.
type BatchReqOpts struct {
	// BatchSize is the internal batch size used per API request; the IPinfo
	// API has a maximum batch size, but the batch request functions available
	// in this library do not. Therefore the library chunks the input slices
	// internally into chunks of size `BatchSize`, clipping to the maximum
	// allowed by the IPinfo API.
	//
	// 0 means to use the default batch size which is the max allowed by the
	// IPinfo API.
	BatchSize uint32

	// TimeoutPerBatch is the timeout in seconds that each batch of size
	// `BatchSize` will have for its own request.
	//
	// 0 means to use a default of 5 seconds; any negative number will turn it
	// off; turning it off does _not_ disable the effects of `TimeoutTotal`.
	TimeoutPerBatch int64

	// TimeoutTotal is the total timeout in seconds for all batch requests in a
	// batch request function to complete.
	//
	// 0 means no total timeout; `TimeoutPerBatch` will still apply.
	TimeoutTotal uint64

	// ConcurrentBatchRequestsLimit is the maximum number of concurrent batch
	// requests that will be mid-flight for inputs that exceed the batch limit.
	//
	// 0 means to use a default of 8; any negative number implies unlimited concurrency.
	ConcurrentBatchRequestsLimit int

	// Filter, if turned on, will filter out a URL whose value was deemed empty
	// on the server.
	Filter bool
}

/* GENERIC */

// GetBatch does a batch request for all `urls` at once.
func GetBatch(
	urls []string,
	opts BatchReqOpts,
) (Batch, error) {
	return DefaultClient.GetBatch(urls, opts)
}

// GetBatch does a batch request for all `urls` at once.
func (c *Client) GetBatch(
	urls []string,
	opts BatchReqOpts,
) (Batch, error) {
	var batchSize int
	var timeoutPerBatch int64
	var maxConcurrentBatchRequests int
	var totalTimeoutCtx context.Context
	var totalTimeoutCancel context.CancelFunc
	var lookupUrls []string
	var result Batch
	var mu sync.Mutex

	// if the cache is available, filter out URLs already cached.
	result = make(Batch, len(urls))
	if c.Cache != nil {
		lookupUrls = make([]string, 0, len(urls)/2)
		for _, url := range urls {
			if res, err := c.Cache.Get(cacheKey(url)); err == nil {
				result[url] = res
			} else {
				lookupUrls = append(lookupUrls, url)
			}
		}
	} else {
		lookupUrls = urls
	}

	// everything cached; exit early.
	if len(lookupUrls) == 0 {
		return result, nil
	}

	// use correct batch size; default/clip to `batchMaxSize`.
	if opts.BatchSize == 0 || opts.BatchSize > batchMaxSize {
		batchSize = batchMaxSize
	} else {
		batchSize = int(opts.BatchSize)
	}

	// use correct concurrent requests limit; either default or user-provided.
	if opts.ConcurrentBatchRequestsLimit == 0 {
		maxConcurrentBatchRequests = batchDefaultConcurrentRequestsLimit
	} else {
		maxConcurrentBatchRequests = opts.ConcurrentBatchRequestsLimit
	}

	// use correct timeout per batch; either default or user-provided.
	if opts.TimeoutPerBatch == 0 {
		timeoutPerBatch = batchReqTimeoutDefault
	} else {
		timeoutPerBatch = opts.TimeoutPerBatch
	}

	// use correct timeout total; either ignore it or apply user-provided.
	if opts.TimeoutTotal > 0 {
		totalTimeoutCtx, totalTimeoutCancel = context.WithTimeout(
			context.Background(),
			time.Duration(opts.TimeoutTotal)*time.Second,
		)
		defer totalTimeoutCancel()
	} else {
		totalTimeoutCtx = context.Background()
	}

	errg, ctx := errgroup.WithContext(totalTimeoutCtx)
	errg.SetLimit(maxConcurrentBatchRequests)
	for i := 0; i < len(lookupUrls); i += batchSize {
		end := i + batchSize
		if end > len(lookupUrls) {
			end = len(lookupUrls)
		}

		urlsChunk := lookupUrls[i:end]
		errg.Go(func() error {
			var postURL string

			// prepare request.

			var timeoutPerBatchCtx context.Context
			var timeoutPerBatchCancel context.CancelFunc
			if timeoutPerBatch > 0 {
				timeoutPerBatchCtx, timeoutPerBatchCancel = context.WithTimeout(
					ctx,
					time.Duration(timeoutPerBatch)*time.Second,
				)
				defer timeoutPerBatchCancel()
			} else {
				timeoutPerBatchCtx = context.Background()
			}

			if opts.Filter {
				postURL = "batch?filter=1"
			} else {
				postURL = "batch"
			}

			jsonArrStr, err := json.Marshal(urlsChunk)
			if err != nil {
				return err
			}
			jsonBuf := bytes.NewBuffer(jsonArrStr)

			req, err := c.newRequest(timeoutPerBatchCtx, "POST", postURL, jsonBuf)
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			// temporarily make a new local result map so that we can read the
			// network data into it; once we have it local we'll merge it with
			// `result` in a concurrency-safe way.
			localResult := new(batch)
			if _, err := c.do(req, localResult); err != nil {
				return err
			}

			// update final result.
			mu.Lock()
			defer mu.Unlock()
			for k, v := range *localResult {
				if strings.HasPrefix(k, "AS") {
					decodedV := new(ASNDetails)
					if err := json.Unmarshal(v, decodedV); err != nil {
						return err
					}

					decodedV.setCountryName()
					result[k] = decodedV
				} else if net.ParseIP(k) != nil {
					decodedV := new(Core)
					if err := json.Unmarshal(v, decodedV); err != nil {
						return err
					}

					decodedV.setCountryName()
					result[k] = decodedV
				} else {
					decodedV := new(interface{})
					if err := json.Unmarshal(v, decodedV); err != nil {
						return err
					}

					result[k] = decodedV
				}
			}

			return nil
		})
	}
	if err := errg.Wait(); err != nil {
		return result, err
	}

	// we delay inserting into the cache until now because:
	// 1. it's likely more cache-line friendly.
	// 2. doing it while updating `result` inside the request workers would be
	//    problematic if the cache is external since we take a mutex lock for
	//    that entire period.
	if c.Cache != nil {
		for _, url := range lookupUrls {
			if v, exists := result[url]; exists {
				if err := c.Cache.Set(cacheKey(url), v); err != nil {
					// NOTE: still return the result even if the cache fails.
					return result, err
				}
			}
		}
	}

	return result, nil
}

/* CORE (net.IP) */

// GetIPInfoBatch does a batch request for all `ips` at once.
func GetIPInfoBatch(
	ips []net.IP,
	opts BatchReqOpts,
) (BatchCore, error) {
	return DefaultClient.GetIPInfoBatch(ips, opts)
}

// GetIPInfoBatch does a batch request for all `ips` at once.
func (c *Client) GetIPInfoBatch(
	ips []net.IP,
	opts BatchReqOpts,
) (BatchCore, error) {
	ipstrs := make([]string, 0, len(ips))
	if c.Token == "" {
		return nil, fmt.Errorf("invalid token")
	}
	for _, ip := range ips {
		if ip == nil {
			continue
		}
		ipstrs = append(ipstrs, ip.String())
	}

	return c.GetIPStrInfoBatch(ipstrs, opts)
}

/* CORE (string) */

// GetIPStrInfoBatch does a batch request for all `ips` at once.
func GetIPStrInfoBatch(
	ips []string,
	opts BatchReqOpts,
) (BatchCore, error) {
	return DefaultClient.GetIPStrInfoBatch(ips, opts)
}

// GetIPStrInfoBatch does a batch request for all `ips` at once.
func (c *Client) GetIPStrInfoBatch(
	ips []string,
	opts BatchReqOpts,
) (BatchCore, error) {
	intermediateRes, err := c.GetBatch(ips, opts)

	// if we have items in the result, don't throw them away; we'll convert
	// below and return the error together if it existed.
	if err != nil && len(intermediateRes) == 0 {
		return nil, err
	}

	res := make(BatchCore, len(intermediateRes))
	for k, v := range intermediateRes {
		res[k] = v.(*Core)
	}

	return res, err
}

/* ASN */

// GetASNDetailsBatch does a batch request for all `asns` at once.
func GetASNDetailsBatch(
	asns []string,
	opts BatchReqOpts,
) (BatchASNDetails, error) {
	return DefaultClient.GetASNDetailsBatch(asns, opts)
}

// GetASNDetailsBatch does a batch request for all `asns` at once.
func (c *Client) GetASNDetailsBatch(
	asns []string,
	opts BatchReqOpts,
) (BatchASNDetails, error) {
	intermediateRes, err := c.GetBatch(asns, opts)

	// if we have items in the result, don't throw them away; we'll convert
	// below and return the error together if it existed.
	if err != nil && len(intermediateRes) == 0 {
		return nil, err
	}

	res := make(BatchASNDetails, len(intermediateRes))
	for k, v := range intermediateRes {
		res[k] = v.(*ASNDetails)
	}
	return res, err
}
