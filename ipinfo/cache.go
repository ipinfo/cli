package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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

type CacheItemString struct {
	LastAccessed time.Time `json:"lat"`
	TimeStamp    time.Time `json:"ts"`
	Data         string    `json:"d"`
}

type CacheItemMap struct {
	LastAccessed time.Time              `json:"lat"`
	TimeStamp    time.Time              `json:"ts"`
	Data         map[string]interface{} `json:"d"`
}

type CacheItemCore struct {
	LastAccessed time.Time    `json:"lat"`
	TimeStamp    time.Time    `json:"ts"`
	Data         *ipinfo.Core `json:"d"`
}

type CacheItemASN struct {
	LastAccessed time.Time          `json:"lat"`
	TimeStamp    time.Time          `json:"ts"`
	Data         *ipinfo.ASNDetails `json:"d"`
}

// Create a new Boltdb-based cache.
func NewBoltdbCache() (*BoltdbCache, error) {
	// get path to database file (<os-dependent-prefix>/ipinfo/cache.boltdb).
	path, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(filepath.Join(path, "ipinfo"), "cache.boltdb")

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

// Sets `key` to `val` in the cache.
//
// This implements `Set` from the IPinfo Go SDK cache interface.
func (c *BoltdbCache) Set(key string, val interface{}) error {
	isFull, err := c.isFull()
	if err != nil {
		return err
	}
	if isFull {
		c.delBulk()
	}
	d, err := c.encode(time.Now(), time.Now(), val)
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

// Gets the value associated with `key` in the cache.
//
// This implements `Get` from the IPinfo Go SDK cache interface.
func (c *BoltdbCache) Get(key string) (interface{}, error) {
	var i interface{}
	var ts time.Time
	err := c.db.View(func(t *bbolt.Tx) error {
		val := t.Bucket(II_CACHE_BUCKET).Get([]byte(key))
		if val == nil {
			return errors.New("key does not exist")
		}

		var err error
		_, ts, i, err = c.decode(val)

		return err
	})
	if err != nil {
		return nil, err
	}
	c.SetLAT(ts, key, i)
	if isOlderThanOneDay(ts) {
		err := c.delCache(key)
		if err != nil {
			return nil, err
		}
	}

	return i, nil
}

// Encodes some data into raw bytes for the cache.
func (c *BoltdbCache) encode(lat time.Time, ts time.Time, d interface{}) ([]byte, error) {
	// get the right output type.
	var t byte
	var i interface{}
	switch dConcrete := d.(type) {
	case string:
		t = II_CACHE_VTYPE_STRING
		i = CacheItemString{LastAccessed: lat, TimeStamp: ts, Data: dConcrete}
	case map[string]interface{}:
		t = II_CACHE_VTYPE_MAP
		i = CacheItemMap{LastAccessed: lat, TimeStamp: ts, Data: dConcrete}
	case *ipinfo.ASNDetails:
		t = II_CACHE_VTYPE_ASN
		i = CacheItemASN{LastAccessed: lat, TimeStamp: ts, Data: dConcrete}
	case *ipinfo.Core:
		t = II_CACHE_VTYPE_CORE
		i = CacheItemCore{LastAccessed: lat, TimeStamp: ts, Data: dConcrete}
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
func (c *BoltdbCache) decode(data []byte) (time.Time, time.Time, interface{}, error) {
	ts := time.Now()

	// last byte contains type info.
	t := data[len(data)-1]

	// get remaining data minus type byte.
	data = data[:len(data)-1]

	// get right decode type.
	switch t {
	case II_CACHE_VTYPE_STRING:
		i := CacheItemString{}
		err := json.Unmarshal(data, &i)
		if err != nil {
			return ts, ts, i, err
		}
		return i.LastAccessed, i.TimeStamp, i.Data, nil
	case II_CACHE_VTYPE_MAP:
		i := CacheItemMap{}
		err := json.Unmarshal(data, &i)
		if err != nil {
			return ts, ts, i, err
		}
		return i.LastAccessed, i.TimeStamp, i.Data, nil
	case II_CACHE_VTYPE_CORE:
		i := CacheItemCore{}
		err := json.Unmarshal(data, &i)
		if err != nil {
			return ts, ts, i, err
		}
		return i.LastAccessed, i.TimeStamp, i.Data, nil
	case II_CACHE_VTYPE_ASN:
		i := CacheItemASN{}
		err := json.Unmarshal(data, &i)
		if err != nil {
			return ts, ts, i, err
		}
		return i.LastAccessed, i.TimeStamp, i.Data, nil
	default:
		return ts, ts, nil, fmt.Errorf("unrecognized type '%v' in cache decoding", t)
	}
}

func isOlderThanOneDay(t time.Time) bool {
	return time.Since(t) > 24*time.Hour
}

func (c *BoltdbCache) delCache(key string) error {
	return c.db.Update(func(t *bbolt.Tx) error {
		err := t.Bucket(II_CACHE_BUCKET).Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("something went wrong while deleting cache: %w", err)
		}
		return nil
	})
}
func (c *BoltdbCache) SetLAT(ttl time.Time, key string, val interface{}) error {
	d, err := c.encode(time.Now(), ttl, val)
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

func (c *BoltdbCache) isFull() (bool, error) {
	var tKey int
	err := c.db.View(func(t *bbolt.Tx) error {
		tKey = t.Bucket(II_CACHE_BUCKET).Stats().KeyN
		return nil
	})
	if err != nil {
		return false, err
	}
	if tKey > 5 {
		return true, nil
	}
	return false, nil
}

type latKeys struct {
	k   string
	lat time.Time
}

type latKeysArr []latKeys

func (s latKeysArr) Len() int {
	return len(s)
}

func (s latKeysArr) Less(i, j int) bool {
	return s[i].lat.Before(s[j].lat)
}

func (s latKeysArr) Swap(i, j int) {
	tmp := s[i]
	s[i] = tmp
	s[j] = tmp
}

func (c *BoltdbCache) delBulk() {
	var latCustom latKeysArr
	c.db.Update(func(t *bbolt.Tx) error {
		err := t.Bucket(II_CACHE_BUCKET).ForEach(func(k, v []byte) error {
			ltime, _, _, _ := c.decode(v)
			latCustom = append(latCustom, latKeys{
				k:   string(k),
				lat: ltime,
			})
			return nil
		})
		if err != nil {
			return err
		}
		sort.Sort(latCustom)
		for i := 0; i < 5; i++ {
			err := t.Bucket(II_CACHE_BUCKET).Delete([]byte(latCustom[i].k))
			if err != nil {
				return err
			}
		}
		return nil
	})

}
