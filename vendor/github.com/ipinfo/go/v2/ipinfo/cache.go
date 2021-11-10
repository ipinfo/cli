package ipinfo

import (
	"fmt"

	"github.com/ipinfo/go/v2/ipinfo/cache"
)

const cacheKeyVsn = "2"

// Cache represents the internal cache used by the IPinfo client.
type Cache struct {
	cache.Interface
}

// NewCache creates a new cache given a specific engine.
func NewCache(engine cache.Interface) *Cache {
	return &Cache{Interface: engine}
}

// return a versioned cache key.
func cacheKey(k string) string {
	return fmt.Sprintf("%s:%s", k, cacheKeyVsn)
}
