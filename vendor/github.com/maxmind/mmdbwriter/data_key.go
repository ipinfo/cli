package mmdbwriter

import (
	"bytes"
	"crypto/sha256"
	"hash"

	"github.com/maxmind/mmdbwriter/mmdbtype"
)

// keyWriter is similar to dataWriter but it will never use pointers. This
// will produce a unique key for the type.
type keyWriter struct {
	*bytes.Buffer
	sha256 hash.Hash
}

func newKeyWriter() *keyWriter {
	return &keyWriter{Buffer: &bytes.Buffer{}, sha256: sha256.New()}
}

// This is just a quick hack. I am sure there is
// something better.
func (kw *keyWriter) key(t mmdbtype.DataType) ([]byte, error) {
	kw.Truncate(0)
	kw.sha256.Reset()
	_, err := t.WriteTo(kw)
	if err != nil {
		return nil, err
	}
	if _, err := kw.WriteTo(kw.sha256); err != nil {
		return nil, err
	}
	return kw.sha256.Sum(nil), nil
}

func (kw *keyWriter) WriteOrWritePointer(t mmdbtype.DataType) (int64, error) {
	return t.WriteTo(kw)
}
