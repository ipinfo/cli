package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

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
	d, err := c.encode(val)
	if err != nil {
		return err
	}

	return c.db.Update(func(t *bbolt.Tx) error {
		err := t.Bucket(II_CACHE_BUCKET).Put([]byte(key), []byte(d))
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
	err := c.db.View(func(t *bbolt.Tx) error {
		val := t.Bucket(II_CACHE_BUCKET).Get([]byte(key))
		if val == nil {
			return errors.New("key does not exist")
		}

		var err error
		i, err = c.decode(val)
		return err
	})
	return i, err
}

// Encodes some data into raw bytes for the cache.
func (c *BoltdbCache) encode(i interface{}) ([]byte, error) {
	// encode as json.
	d, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("could not encode bytes: %w", err)
	}

	// add a single type byte at the end for decoding purposes.
	var t byte
	switch i.(type) {
	case string:
		t = II_CACHE_VTYPE_STRING
	case map[string]interface{}:
		t = II_CACHE_VTYPE_MAP
	case *ipinfo.ASNDetails:
		t = II_CACHE_VTYPE_ASN
	case *ipinfo.Core:
		t = II_CACHE_VTYPE_CORE
	default:
		return nil, fmt.Errorf("unrecognized type '%v' in cache encoding", reflect.TypeOf(i))
	}
	d = append(d, t)

	return d, nil

}

// function to decode rawdata and return
func (c *BoltdbCache) decode(data []byte) (interface{}, error) {
	// last byte contains type info.
	t := data[len(data)-1]

	// get remaining data minus type byte.
	data = data[:len(data)-1]

	// decode according to type.
	switch t {
	case II_CACHE_VTYPE_STRING:
		var d string
		err := json.Unmarshal(data, &d)
		return d, err
	case II_CACHE_VTYPE_MAP:
		var d map[string]interface{}
		err := json.Unmarshal(data, &d)
		return d, err
	case II_CACHE_VTYPE_CORE:
		var d *ipinfo.Core
		err := json.Unmarshal(data, &d)
		return d, err
	case II_CACHE_VTYPE_ASN:
		var d *ipinfo.ASNDetails
		err := json.Unmarshal(data, &d)
		return d, err
	default:
		return nil, fmt.Errorf("unrecognized type '%v' in cache decoding", t)
	}
}
