package ipUtils_test

import (
	"testing"

	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/stretchr/testify/assert"
)

func TestIPRangeStrFromStr(t *testing.T) {
	var r ipUtils.IPRangeStr
	var err error

	r, err = ipUtils.IPRangeStrFromStr("-")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr(",")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1,")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1-")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1,")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1-")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1-")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,,")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,-")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1--")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,1.1.1")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,1.1.1,")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1-1.1.1,")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1-1.1.1-")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,1.1.1.2,")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,1.1.1.2-")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,1.1.1.2.3")
	assert.Equal(t, err, ipUtils.ErrNotIPRange)

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1-1.1.1.2")
	assert.Nil(t, err)
	assert.Equal(t, r.Start, "1.1.1.1")
	assert.Equal(t, r.End, "1.1.1.2")

	r, err = ipUtils.IPRangeStrFromStr("1.1.1.1,1.1.1.2")
	assert.Nil(t, err)
	assert.Equal(t, r.Start, "1.1.1.1")
	assert.Equal(t, r.End, "1.1.1.2")
}
