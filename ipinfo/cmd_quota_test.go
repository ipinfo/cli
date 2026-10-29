package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintStats(t *testing.T) {
	tests := []struct {
		name     string
		requests QuotaRequests
	}{
		{"under limit", QuotaRequests{Day: 100, Month: 1000, Limit: 50000, Remaining: 49000}},
		{"at limit", QuotaRequests{Day: 100, Month: 50000, Limit: 50000, Remaining: 0}},
		{"over limit", QuotaRequests{Day: 47038, Month: 50282, Limit: 50000, Remaining: 0}},
		// a response without a limit must not divide by zero.
		{"zero limit", QuotaRequests{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quota := &QuotaBody{Requests: tt.requests}
			assert.NotPanics(t, func() { printStats(false, quota) })
		})
	}
}

func TestGetUsageBar(t *testing.T) {
	assert.Equal(t, "["+strings.Repeat("#", 50)+strings.Repeat(" ", 50)+"] 50%\n", getUsageBar(50))
	assert.Equal(t, "["+strings.Repeat("#", 100)+"] 101%\n", getUsageBar(101))
}
