package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipinfo/go/v2/ipinfo"
	bolt "go.etcd.io/bbolt"
)

const (
	string_val = iota
	map_val
	core_val
	ANSDetails_val
	other_val
)

type db_cache struct {
	db *bolt.DB
}

//function to setup boltdb
func setupDB() (*db_cache, error) {
	config_path, err := os.UserConfigDir()
	db_path := filepath.Join(config_path, "ipinfo_cache")
	if err != nil {
		return nil, err
	}
	//opening database file will create if doesn't exist already
	db, err := bolt.Open(db_path, 0660, nil)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	//creating root bucket
	err = db.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists([]byte("ipinfo_cache"))
		if err != nil {
			return fmt.Errorf("error in creating root bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error setting up the database: %v", err)
	}
	return &db_cache{
		db: db,
	}, nil

}

//implementation of Set function from interface
func (db *db_cache) Set(key string, value interface{}) error {
	raw_data, err := encodeBytes(value)
	if err != nil {
		return err
	}
	err = db.db.Update(func(t *bolt.Tx) error {
		err := t.Bucket([]byte("ipinfo_cache")).Put([]byte(key), []byte(raw_data))
		if err != nil {
			return fmt.Errorf("error in adding data: %v", err)
		}
		return err

	})
	return err

}

//implemnentation of Get function from interface
func (db *db_cache) Get(key string) (interface{}, error) {
	var value []byte
	err := db.db.Update(func(t *bolt.Tx) error {
		value = t.Bucket([]byte("ipinfo_cache")).Get([]byte(key))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("something went wrong in fetching data")

	}

	return decodeByte(value), nil

}

//function to encode interface to raw bytes
func encodeBytes(i interface{}) ([]byte, error) {
	output, err := json.Marshal(i)
	switch i.(type) {
	case string:
		output = append(output, string_val)
	case map[string]interface{}:
		output = append(output, map_val)
	case ipinfo.ASNDetails:
		output = append(output, ANSDetails_val)
	case ipinfo.Core:
		output = append(output, core_val)
	default:
		output = append(output, other_val)
	}

	if err != nil {
		return nil, fmt.Errorf("error in encodeerting to bytes")
	}
	return output, nil

}

// function to decode rawdata and return
func decodeByte(data []byte) interface{} {

	last := data[len(data)-1]
	data = data[:len(data)-1]
	switch last {
	case string_val:
		var st_data string
		json.Unmarshal(data, &st_data)
		return st_data
	case map_val:
		var map_data map[string]interface{}
		json.Unmarshal(data, &map_data)
		return map_data
	case core_val:
		var core_data ipinfo.Core
		json.Unmarshal(data, &core_data)
		return core_data
	case ANSDetails_val:
		var ASN_data ipinfo.ASNDetails
		json.Unmarshal(data, &ASN_data)
		return ASN_data
	default:
		return nil

	}
}
