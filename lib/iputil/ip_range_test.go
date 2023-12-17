package iputil_test

import (
	"testing"

	"github.com/ipinfo/cli/lib/iputil"
	"github.com/stretchr/testify/assert"
)

func TestIPRangeStrFromStr(t *testing.T) {
	var r iputil.IPRangeStr
	var err error

	r, err = iputil.IPRangeStrFromStr("-")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr(",")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1,")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1-")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1,")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1-")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1-")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,,")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,-")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1--")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,1.1.1")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,1.1.1,")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1-1.1.1,")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1-1.1.1-")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,1.1.1.2,")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,1.1.1.2-")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,1.1.1.2.3")
	assert.Equal(t, err, iputil.ErrNotIPRange)

	r, err = iputil.IPRangeStrFromStr("1.1.1.1-1.1.1.2")
	assert.Nil(t, err)
	assert.Equal(t, r.Start, "1.1.1.1")
	assert.Equal(t, r.End, "1.1.1.2")

	r, err = iputil.IPRangeStrFromStr("1.1.1.1,1.1.1.2")
	assert.Nil(t, err)
	assert.Equal(t, r.Start, "1.1.1.1")
	assert.Equal(t, r.End, "1.1.1.2")
}
