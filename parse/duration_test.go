package parse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDuration(t *testing.T) {
	tests := map[string]time.Duration{
		"PT10S":     10 * time.Second,
		"PT100S":    100 * time.Second,
		"PT8M":      8 * time.Minute,
		"PT5M40S":   5*time.Minute + 40*time.Second,
		"PT4H":      4 * time.Hour,
		"PT8H2M":    8*time.Hour + 2*time.Minute,
		"PT2H10M2S": 2*time.Hour + 10*time.Minute + 2*time.Second,
		"P4W":       4 * 7 * 24 * time.Hour,
		"P7D":       7 * 24 * time.Hour,
		"P7DT4H10S": 7*24*time.Hour + 4*time.Hour + 10*time.Second,
	}

	for raw, expected := range tests {
		t.Run(raw, func(t *testing.T) {
			t.Run("no sign", testParseDuration(raw, expected))
			t.Run("plus sign", testParseDuration("+"+raw, +expected))
			t.Run("minus sign", testParseDuration("-"+raw, -expected))
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
