package ipinfo

import (
	"github.com/ipinfo/go/v2/ipinfo/cache"
)

// Cache represents the internal cache used by the IPinfo client.
type Cache struct {
	cache.Interface
}

// NewCache creates a new cache given a specific engine.
func NewCache(engine cache.Interface) *Cache {
	return &Cache{Interface: engine}
}
