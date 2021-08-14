package lib_test

import (
	"testing"

	"github.com/ipinfo/cli/lib"
	"github.com/stretchr/testify/assert"
)

func TestIPRangeStrFromStr(t *testing.T) {
	var r IPRange
	var err error

	r, err = lib.IPRangeStrFromStr("-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr(",")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1--")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,1.1.1")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1-1.1.1,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1-1.1.1-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,1.1.1.2,")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,1.1.1.2-")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1,1.1.1.2.3")
	assert.Equal(t, err, lib.ErrNotIPRange)

	r, err = lib.IPRangeStrFromStr("1.1.1.1-1.1.1.2")
	assert.Nil(t, err)
	assert.Equal(t, r.Start, "1.1.1.1")
	assert.Equal(t, r.End, "1.1.1.2")

	r, err = lib.IPRangeStrFromStr("1.1.1.1,1.1.1.2")
	assert.Nil(t, err)
	assert.Equal(t, r.Start, "1.1.1.1")
	assert.Equal(t, r.End, "1.1.1.2")
}
