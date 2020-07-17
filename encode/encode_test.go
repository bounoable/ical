package encode_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bounoable/ical/encode"
	"github.com/bounoable/ical/internal/testutil"
	"github.com/bounoable/ical/parse"
	"github.com/stretchr/testify/assert"
)

func TestCalendar(t *testing.T) {
	tests := []struct {
		calendar parse.Calendar
		expected string
	}{
		{
			calendar: parse.Calendar{
				Properties: []parse.Property{
					testutil.Property("X-FOO", "bar", parse.Parameters{
						"foo": []string{"bar", "baz"},
					}),
					testutil.Property("PRODID", "123456abcdef", nil),
					testutil.Property("VERSION", "2.0", nil),
					testutil.Property("CALSCALE", "GREGORIAN", nil),
					testutil.Property("METHOD", "REQUEST", nil),
				},
				Events: []parse.Event{
					{
						Properties: []parse.Property{
							testutil.Property("UID", "111111111111", nil),
							testutil.Property("DTSTART", "20200101", nil),
							testutil.Property("DTEND", "20200301T103000", parse.Parameters{
								"VALUE": []string{"DATE-TIME"},
							}),
							testutil.Property("SUMMARY", "foo summary", nil),
							testutil.Property("DESCRIPTION", "this is a long description that should be folded onto the next line", nil),
						},
						// UID:         "111111111111",
						// Summary:     "foo summary",
						// Description: "very long description that should be folded onto the next line",
						// Start:       time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
						// End:         time.Date(2020, time.March, 1, 10, 30, 0, 0, time.Local),
					},
				},
			},
			expected: `BEGIN:VCALENDAR
X-FOO;foo=bar,baz:bar
PRODID:123456abcdef
VERSION:2.0
CALSCALE:GREGORIAN
METHOD:REQUEST
BEGIN:VEVENT
UID:111111111111
DTSTART:20200101
DTEND;VALUE=DATE-TIME:20200301T103000
SUMMARY:foo summary
DESCRIPTION:this is a long description that should be folded onto the next 
 line
END:VEVENT
END:VCALENDAR`,
		},
	}

	for i, test := range tests {
		test.expected = strings.ReplaceAll(test.expected, "\n", "\r\n")
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var buf strings.Builder
			err := encode.Calendar(test.calendar, &buf)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expected, buf.String())
		})
	}
}
