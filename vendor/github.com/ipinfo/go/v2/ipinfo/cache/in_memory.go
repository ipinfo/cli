package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const defaultExpiration = 24 * time.Hour

// InMemory is an implementation of the cache interface which stores values
// in-memory.
type InMemory struct {
	cache      *cache.Cache
	expiration time.Duration
}

// NewInMemory creates a new InMemory instance with default values.
func NewInMemory() *InMemory {
	return &InMemory{
		cache:      cache.New(-1, defaultExpiration),
		expiration: defaultExpiration,
	}
}

// WithExpiration updates the expiration value of `c`.
func (c *InMemory) WithExpiration(d time.Duration) *InMemory {
	c.expiration = d
	return c
}

// Get retrieves a value from the InMemory cache implementation.
func (c *InMemory) Get(key string) (interface{}, error) {
	v, found := c.cache.Get(key)
	if !found {
		return nil, ErrNotFound
	}
	return v, nil
}

// Set sets a value for a key in the InMemory cache implementation.
func (c *InMemory) Set(key string, value interface{}) error {
	c.cache.Set(key, value, c.expiration)
	return nil
}

// Check if InMemory implements cache.Interface
var _ Interface = (*InMemory)(nil)
