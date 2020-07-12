package parse_test

import (
	"testing"
	"time"

	"github.com/bounoable/ical/lex"
	"github.com/bounoable/ical/parse"
	"github.com/stretchr/testify/assert"
)

func TestParse_event(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected parse.Event
	}{
		{
			name: "DTSTART - DTEND (UTC)",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTAMP:19970901T130000Z
DTSTART:19970903T163000Z
DTEND:19970903T190000Z
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("DTSTAMP", "19970901T130000Z", nil),
					property("DTSTART", "19970903T163000Z", nil),
					property("DTEND", "19970903T190000Z", nil),
				},
				Timestamp: time.Date(1997, time.September, 1, 13, 0, 0, 0, time.UTC),
				Start:     time.Date(1997, time.September, 3, 16, 30, 0, 0, time.UTC),
				End:       time.Date(1997, time.September, 3, 19, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "annual event",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTAMP:19970901T130000Z
DTSTART;VALUE=DATE:19971102
RRULE:FREQ=YEARLY
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("DTSTAMP", "19970901T130000Z", nil),
					property("DTSTART", "19971102", parse.Parameters{
						"VALUE": []string{"DATE"},
					}),
					property("RRULE", "FREQ=YEARLY", nil),
				},
				Timestamp: time.Date(1997, time.September, 1, 13, 0, 0, 0, time.UTC),
				Start:     time.Date(1997, time.November, 2, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "multi day event",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTAMP:20070423T123432Z
DTSTART;VALUE=DATE:20070628
DTEND;VALUE=DATE:20070709
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("DTSTAMP", "20070423T123432Z", nil),
					property("DTSTART", "20070628", parse.Parameters{
						"VALUE": []string{"DATE"},
					}),
					property("DTEND", "20070709", parse.Parameters{
						"VALUE": []string{"DATE"},
					}),
				},
				Timestamp: time.Date(2007, time.April, 23, 12, 34, 32, 0, time.UTC),
				Start:     time.Date(2007, time.June, 28, 0, 0, 0, 0, time.Local),
				End:       time.Date(2007, time.July, 9, 0, 0, 0, 0, time.Local),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			items, err := lex.Text(test.input)
			if err != nil {
				t.Fatal(err)
			}

			cal, err := parse.Items(items)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expected, cal.Events[0])
		})
	}
}
