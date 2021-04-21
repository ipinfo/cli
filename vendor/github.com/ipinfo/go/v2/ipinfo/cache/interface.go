package cache

import "errors"

var (
	// ErrNotFound means that the key was not found.
	ErrNotFound = errors.New("key not found")
)

// Interface is the cache interface that all cache implementations must adhere
// to at the minimum to be usable in the IPinfo client.
//
// Note that all implementations must be concurrency-safe.
type Interface interface {
	// Get a value from the cache given a key.
	//
	// This must be concurrency-safe.
	Get(key string) (interface{}, error)

	// Set a key to value mapping in the cache.
	//
	// This must be concurrency-safe.
	Set(key string, value interface{}) error
}
