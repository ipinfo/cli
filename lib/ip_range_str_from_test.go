package lib_test

import (
	"testing"

	"github.com/ipinfo/cli/lib"
	"github.com/stretchr/testify/assert"
)

func TestIPRangeStrFromRangeStr(t *testing.T) {
	var s string // start
	var e string // end
	var err error

	s, e, err = lib.IPRangeStrFromRangeStr("-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr(",")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1--")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,1.1.1")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1-1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1-1.1.1-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,1.1.1.2,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,1.1.1.2-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,1.1.1.2.3")
	assert.Equal(t, err, lib.ErrNotIPRange)

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1-1.1.1.2")
	assert.Nil(t, err)
	assert.Equal(t, s, "1.1.1.1")
	assert.Equal(t, e, "1.1.1.2")

	s, e, err = lib.IPRangeStrFromRangeStr("1.1.1.1,1.1.1.2")
	assert.Nil(t, err)
	assert.Equal(t, s, "1.1.1.1")
	assert.Equal(t, e, "1.1.1.2")
}
