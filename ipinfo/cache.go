package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"time"

	"github.com/ipinfo/go/v2/ipinfo"
	"go.etcd.io/bbolt"
)

const (
	II_CACHE_VTYPE_STRING = iota
	II_CACHE_VTYPE_MAP
	II_CACHE_VTYPE_CORE
	II_CACHE_VTYPE_ASN
)

var II_CACHE_BUCKET []byte = []byte("ii")

type BoltdbCache struct {
	db *bbolt.DB
}

type CacheItem struct {
	LastAccessed time.Time `json:"lat"`
	Created      time.Time `json:"cr"`
}

type CacheItemString struct {
	CacheItem
	Data string `json:"d"`
}

type CacheItemMap struct {
	CacheItem
	Data map[string]interface{} `json:"d"`
}

type CacheItemCore struct {
	CacheItem
	Data *ipinfo.Core `json:"d"`
}

type CacheItemASN struct {
	CacheItem
	Data *ipinfo.ASNDetails `json:"d"`
}

// Returns the path to the cache database file.
func BoltdbCachePath() (string, error) {
	confDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(confDir, "cache.boltdb"), nil
}

// Create a new Boltdb-based cache.
func NewBoltdbCache() (*BoltdbCache, error) {
	// get path to database file.
	path, err := BoltdbCachePath()
	if err != nil {
		return nil, err
	}

	// open db.
	db, err := bbolt.Open(path, 0660, nil)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// creating root bucket.
	err = db.Update(func(t *bbolt.Tx) error {
		_, err := t.CreateBucketIfNotExists(II_CACHE_BUCKET)
		if err != nil {
			return fmt.Errorf("error creating root db bucket: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error setting up db: %w", err)
	}

	return &BoltdbCache{db: db}, nil
}

// Gets the value associated with `key` in the cache.
//
// This implements `Get` from the IPinfo Go SDK cache interface.
func (c *BoltdbCache) Get(key string) (interface{}, error) {
	var i interface{}
	var created time.Time
	err := c.db.View(func(t *bbolt.Tx) error {
		val := t.Bucket(II_CACHE_BUCKET).Get([]byte(key))
		if val == nil {
			return errors.New("key does not exist")
		}

		var err error
		_, created, i, err = c.decode(val)
		return err
	})
	if err != nil {
		return nil, err
	}

	// update val.
	if time.Since(created) > 24*time.Hour {
		// has it expired? if so, delete it and return as if we don't have it.
		err := c.del(key)
		if err == nil {
			err = errors.New("key does not exist")
		}
		return nil, err
	}

	// update the same item with an updated last-accessed.
	c.set(key, CacheItem{LastAccessed: time.Now(), Created: created}, i)
	return i, nil
}

// Sets `key` to `val` in the cache.
//
// This implements `Set` from the IPinfo Go SDK cache interface.
func (c *BoltdbCache) Set(key string, val interface{}) error {
	// if the cache is too full, delete a bunch of stuff in bulk before
	// proceeding.
	if c.isTooFull() {
		c.delBulk()
	}

	// set value.
	now := time.Now()
	return c.set(key, CacheItem{LastAccessed: now, Created: now}, val)
}

// Encodes some data into raw bytes for the cache.
func (c *BoltdbCache) encode(
	lat time.Time,
	created time.Time,
	d interface{},
) ([]byte, error) {
	var t byte
	var i interface{}

	// get the right output type.
	cacheItem := CacheItem{
		LastAccessed: lat,
		Created:      created,
	}
	switch dConcrete := d.(type) {
	case string:
		t = II_CACHE_VTYPE_STRING
		i = CacheItemString{CacheItem: cacheItem, Data: dConcrete}
	case map[string]interface{}:
		t = II_CACHE_VTYPE_MAP
		i = CacheItemMap{CacheItem: cacheItem, Data: dConcrete}
	case *ipinfo.ASNDetails:
		t = II_CACHE_VTYPE_ASN
		i = CacheItemASN{CacheItem: cacheItem, Data: dConcrete}
	case *ipinfo.Core:
		t = II_CACHE_VTYPE_CORE
		i = CacheItemCore{CacheItem: cacheItem, Data: dConcrete}
	default:
		return nil, fmt.Errorf("unrecognized type '%v' in cache encoding", reflect.TypeOf(d))
	}

	// encode as json.
	dEncoded, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("could not encode bytes: %w", err)
	}

	// add a single type byte at the end for decoding purposes.
	dEncoded = append(dEncoded, t)

	return dEncoded, nil
}

// function to decode rawdata and return
func (c *BoltdbCache) decode(
	data []byte,
) (time.Time, time.Time, interface{}, error) {
	// last byte contains type info.
	t := data[len(data)-1]

	// get remaining data minus type byte.
	data = data[:len(data)-1]

	// get right decode type.
	var err error
	switch t {
	case II_CACHE_VTYPE_STRING:
		i := CacheItemString{}
		err = json.Unmarshal(data, &i)
		if err == nil {
			return i.LastAccessed, i.Created, i.Data, nil
		}
	case II_CACHE_VTYPE_MAP:
		i := CacheItemMap{}
		err := json.Unmarshal(data, &i)
		if err == nil {
			return i.LastAccessed, i.Created, i.Data, nil
		}
	case II_CACHE_VTYPE_CORE:
		i := CacheItemCore{}
		err := json.Unmarshal(data, &i)
		if err == nil {
			return i.LastAccessed, i.Created, i.Data, nil
		}
	case II_CACHE_VTYPE_ASN:
		i := CacheItemASN{}
		err := json.Unmarshal(data, &i)
		if err == nil {
			return i.LastAccessed, i.Created, i.Data, nil
		}
	default:
		err = fmt.Errorf("unrecognized type '%v' in cache decoding", t)
	}

	return time.Now(), time.Now(), nil, err
}

func (c *BoltdbCache) del(key string) error {
	return c.db.Update(func(t *bbolt.Tx) error {
		err := t.Bucket(II_CACHE_BUCKET).Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("something went wrong while deleting cache: %w", err)
		}
		return nil
	})
}

func (c *BoltdbCache) set(
	key string,
	item CacheItem,
	val interface{},
) error {
	d, err := c.encode(item.LastAccessed, item.Created, val)
	if err != nil {
		return err
	}

	return c.db.Update(func(t *bbolt.Tx) error {
		err := t.Bucket(II_CACHE_BUCKET).Put([]byte(key), d)
		if err != nil {
			return fmt.Errorf("error in adding data: %w", err)
		}

		return nil
	})
}

func (c *BoltdbCache) isTooFull() bool {
	var dbsize int64
	c.db.View(func(t *bbolt.Tx) error {
		dbsize = t.Size()
		return nil
	})
	return dbsize > 1024*1024*1024 // 1 GB
}

type latKey struct {
	k   string
	lat time.Time
}

type latKeys []latKey

func (s latKeys) Len() int {
	return len(s)
}

func (s latKeys) Less(i, j int) bool {
	return s[i].lat.Before(s[j].lat)
}

func (s latKeys) Swap(i, j int) {
	tmp := s[i]
	s[i] = tmp
	s[j] = tmp
}

func (c *BoltdbCache) delBulk() {
	c.db.Update(func(t *bbolt.Tx) error {
		bucket := t.Bucket(II_CACHE_BUCKET)

		// collect all keys with their last-accessed timestamp.
		lats := make(latKeys, 0, bucket.Stats().KeyN)
		err := bucket.ForEach(func(k, v []byte) error {
			lat, _, _, err := c.decode(v)
			if err != nil {
				return err
			}
			lats = append(lats, latKey{k: string(k), lat: lat})
			return nil
		})
		if err != nil {
			return err
		}

		// sort by last-accessed timestamp, and delete the oldest 10k in bulk.
		sort.Sort(lats)
		for i := 0; i < 10000; i++ {
			if err := bucket.Delete([]byte(lats[i].k)); err != nil {
				return err
			}
		}
		return nil
	})
}
