package mmdbwriter

import (
	"bytes"
	"crypto/sha256"
	"hash"

	"github.com/maxmind/mmdbwriter/mmdbtype"
)

// KeyGenerator generates a unique key for record values being inserted into the
// Tree. This is used for deduplicating the values in memory. The default
// KeyGenerator will serialize and hash the whole datastructure. This handles
// the general case well but may be inefficient given the particulars of the
// data.
//
// Please be certain that any key you generate is unique. If there is a
// collision with two different values having the same key, one of the
// values will be overwritten.
//
// The returned byte slice is not stored. You may use the same backing
// array between calls.
type KeyGenerator interface {
	Key(mmdbtype.DataType) ([]byte, error)
}

var _ KeyGenerator = &keyWriter{}

// keyWriter is similar to dataWriter but it will never use pointers. This
// will produce a unique key for the type.
type keyWriter struct {
	*bytes.Buffer
	sha256 hash.Hash
	key    [sha256.Size]byte
}

func newKeyWriter() *keyWriter {
	return &keyWriter{Buffer: &bytes.Buffer{}, sha256: sha256.New()}
}

// This is just a quick hack. I am sure there is
// something better.
func (kw *keyWriter) Key(v mmdbtype.DataType) ([]byte, error) {
	kw.Truncate(0)
	kw.sha256.Reset()
	_, err := v.WriteTo(kw)
	if err != nil {
		return nil, err
	}
	if _, err := kw.WriteTo(kw.sha256); err != nil {
		return nil, err
	}
	return kw.sha256.Sum(kw.key[:0]), nil
}

func (kw *keyWriter) WriteOrWritePointer(t mmdbtype.DataType) (int64, error) {
	return t.WriteTo(kw)
}
