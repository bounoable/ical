package parse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		raw      string
		expected time.Duration
	}{
		{
			raw:      "P10S",
			expected: 10 * time.Second,
		},
		{
			raw:      "P100S",
			expected: 100 * time.Second,
		},
		{
			raw:      "P8M",
			expected: 8 * time.Minute,
		},
		{
			raw:      "P5M40S",
			expected: 5*time.Minute + 40*time.Second,
		},
		{
			raw:      "P4H",
			expected: 4 * time.Hour,
		},
		{
			raw:      "P8H2M",
			expected: 8*time.Hour + 2*time.Minute,
		},
		{
			raw:      "P2H10M2S",
			expected: 2*time.Hour + 10*time.Minute + 2*time.Second,
		},
		{
			raw:      "P4W",
			expected: 4 * 7 * 24 * time.Hour,
		},
		{
			raw:      "P2W7D",
			expected: 2*7*24*time.Hour + 7*24*time.Hour,
		},
	}

	for _, test := range tests {
		t.Run(test.raw, func(t *testing.T) {
			t.Run("no sign", testParseDuration(test.raw, test.expected))
			t.Run("plus sign", testParseDuration("+"+test.raw, +test.expected))
			t.Run("minus sign", testParseDuration("-"+test.raw, -test.expected))
		})
	}
}

func testParseDuration(raw string, expected time.Duration) func(*testing.T) {
	return func(t *testing.T) {
		dur, err := parseDuration(raw)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expected, dur)
	}
}
