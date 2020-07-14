package parse_test

import (
	"context"
	"testing"
	"time"

	"github.com/bounoable/ical/lex"
	"github.com/bounoable/ical/parse"
	"github.com/stretchr/testify/assert"
)

func TestItems(t *testing.T) {
	items := []lex.Item{
		beginCalendar(),
		item(lex.Name, "VERSION"),
		item(lex.Value, "2.0"),
		item(lex.Name, "METHOD"),
		item(lex.Value, "REQUEST"),
		item(lex.Name, "PRODID"),
		item(lex.Value, "-//Example//Product//ID//EN"),
		beginEvent(),
		item(lex.Name, "UID"),
		item(lex.Value, "111111111111"),
		item(lex.Name, "DTSTART"),
		item(lex.ParamName, "VALUE"),
		item(lex.ParamValue, "DATE"),
		item(lex.Value, "20200101"),
		item(lex.Name, "DTEND"),
		item(lex.ParamName, "VALUE"),
		item(lex.ParamValue, "DATE"),
		item(lex.Value, "20200110"),
		endEvent(),
		beginEvent(),
		item(lex.Name, "UID"),
		item(lex.Value, "222222222222"),
		item(lex.Name, "DTSTART"),
		item(lex.ParamName, "VALUE"),
		item(lex.ParamValue, "DATE"),
		item(lex.Value, "20200201"),
		item(lex.Name, "DTEND"),
		item(lex.ParamName, "VALUE"),
		item(lex.ParamValue, "DATE"),
		item(lex.Value, "20200210"),
		item(lex.Name, "DTSTAMP"),
		item(lex.ParamName, "VALUE"),
		item(lex.ParamValue, "DATE-TIME"),
		item(lex.Value, "20200210T103000Z"),
		endEvent(),
		endCalendar(),
	}
	expected := parse.Calendar{
		Properties: []parse.Property{
			property("VERSION", "2.0", nil),
			property("METHOD", "REQUEST", nil),
			property("PRODID", "-//Example//Product//ID//EN", nil),
		},
		Version:   "2.0",
		Method:    "REQUEST",
		ProductID: "-//Example//Product//ID//EN",
		Calscale:  "GREGORIAN",
		Events: []parse.Event{
			{
				Properties: []parse.Property{
					property("UID", "111111111111", nil),
					property("DTSTART", "20200101", parse.Parameters{
						"VALUE": []string{"DATE"},
					}),
					property("DTEND", "20200110", parse.Parameters{
						"VALUE": []string{"DATE"},
					}),
				},
				UID:   "111111111111",
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2020, time.January, 10, 0, 0, 0, 0, time.Local),
			},
			{
				Properties: []parse.Property{
					property("UID", "222222222222", nil),
					property("DTSTART", "20200201", parse.Parameters{
						"VALUE": []string{"DATE"},
					}),
					property("DTEND", "20200210", parse.Parameters{
						"VALUE": []string{"DATE"},
					}),
					property("DTSTAMP", "20200210T103000Z", parse.Parameters{
						"VALUE": []string{"DATE-TIME"},
					}),
				},
				UID:       "222222222222",
				Start:     time.Date(2020, time.February, 1, 0, 0, 0, 0, time.Local),
				End:       time.Date(2020, time.February, 10, 0, 0, 0, 0, time.Local),
				Timestamp: time.Date(2020, time.February, 10, 10, 30, 00, 00, time.UTC),
			},
		},
	}

	res, err := parse.Items(lexItems(items...))
	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}

func TestItems_timeParsing(t *testing.T) {
	tests := map[string]struct {
		items  []lex.Item
		expect func(*testing.T, parse.Calendar)
	}{
		"DATE (default)": {
			items: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "DTSTAMP"),
				item(lex.Value, "20200101"),
				endEvent(),
				endCalendar(),
			},
			expect: func(t *testing.T, cal parse.Calendar) {
				assert.Equal(t, time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local).Unix(), cal.Events[0].Timestamp.Unix())
			},
		},
		"DATE-TIME (local)": {
			items: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE-TIME"),
				item(lex.Value, "20200101T103020"),
				endEvent(),
				endCalendar(),
			},
			expect: func(t *testing.T, cal parse.Calendar) {
				assert.Equal(t, time.Date(2020, time.January, 1, 10, 30, 20, 0, time.Local).Unix(), cal.Events[0].Timestamp.Unix())
			},
		},
		"DATE-TIME (UTC)": {
			items: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE-TIME"),
				item(lex.Value, "20200101T103020Z"),
				endEvent(),
				endCalendar(),
			},
			expect: func(t *testing.T, cal parse.Calendar) {
				assert.Equal(t, time.Date(2020, time.January, 1, 10, 30, 20, 0, time.UTC).Unix(), cal.Events[0].Timestamp.Unix())
			},
		},
		"DATE (malformed as DATE-TIME (local))": {
			items: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20200101T103020"),
				endEvent(),
				endCalendar(),
			},
			expect: func(t *testing.T, cal parse.Calendar) {
				assert.Equal(t, time.Date(2020, time.January, 1, 10, 30, 20, 0, time.Local).Unix(), cal.Events[0].Timestamp.Unix())
			},
		},
		"DATE (malformed as DATE-TIME (UTC))": {
			items: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20200101T103020Z"),
				endEvent(),
				endCalendar(),
			},
			expect: func(t *testing.T, cal parse.Calendar) {
				assert.Equal(t, time.Date(2020, time.January, 1, 10, 30, 20, 0, time.UTC).Unix(), cal.Events[0].Timestamp.Unix())
			},
		},
		"with TZID param": {
			items: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.ParamName, "TZID"),
				item(lex.ParamValue, "America/New_York"),
				item(lex.Value, "20200101"),
				endEvent(),
				endCalendar(),
			},
			expect: func(t *testing.T, cal parse.Calendar) {
				loc, err := time.LoadLocation("America/New_York")
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, time.Date(2020, time.January, 1, 0, 0, 0, 0, loc).Unix(), cal.Events[0].Timestamp.Unix())
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cal, err := parse.Items(lexItems(test.items...))
			assert.Nil(t, err)
			test.expect(t, cal)
		})
	}
}

func TestItems_paramValues(t *testing.T) {
	tests := map[string]struct {
		items  []lex.Item
		expect func(*testing.T, parse.Calendar)
	}{
		"normal string": {
			items: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "X-CUSTOM"),
				item(lex.ParamName, "X-PARAM"),
				item(lex.ParamValue, "foo bar"),
				item(lex.ParamValue, "foo bar baz"),
				item(lex.Value, "bar foo"),
				endEvent(),
				endCalendar(),
			},
			expect: func(t *testing.T, cal parse.Calendar) {
				assert.Equal(t, []string{"foo bar", "foo bar baz"}, cal.Events[0].Properties[0].Params["X-PARAM"])
			},
		},
		"quoted string": {
			items: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "X-CUSTOM"),
				item(lex.ParamName, "X-PARAM"),
				item(lex.ParamValue, "foo bar"),
				item(lex.ParamValue, `"foo bar baz"`),
				item(lex.Value, "bar foo"),
				endEvent(),
				endCalendar(),
			},
			expect: func(t *testing.T, cal parse.Calendar) {
				assert.Equal(t, []string{"foo bar", `"foo bar baz"`}, cal.Events[0].Properties[0].Params["X-PARAM"])
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			cal, err := parse.Items(lexItems(test.items...))
			if err != nil {
				t.Fatal(err)
			}

			test.expect(t, cal)
		})
	}
}

func TestItems_location(t *testing.T) {
	locs := [...]*time.Location{
		time.UTC,
		time.Local,
		loadLocation("America/New_York"),
		loadLocation("Europe/Berlin"),
	}

	for _, loc := range locs {
		t.Run(loc.String(), func(t *testing.T) {
			t.Run("valid layout", parseLocationTest(
				loc,
				"20200101T103000",
				time.Date(2020, time.January, 1, 10, 30, 0, 0, loc),
			))

			t.Run("utc layout", parseLocationTest(
				loc,
				"20200101T103000Z",
				time.Date(2020, time.January, 1, 10, 30, 0, 0, time.UTC),
			))
		})
	}
}

func parseLocationTest(loc *time.Location, layout string, expected time.Time) func(t *testing.T) {
	return func(t *testing.T) {
		items := lexItems(
			beginCalendar(),
			beginEvent(),
			item(lex.Name, "DTSTAMP"),
			item(lex.ParamName, "VALUE"),
			item(lex.ParamValue, "DATE-TIME"),
			item(lex.Value, layout),
			endEvent(),
			endCalendar(),
		)

		cal, err := parse.Items(items, parse.Location(loc))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, cal.Events[0].Timestamp)
	}
}

func TestItems_context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := parse.Items(lexItems(), parse.Context(ctx))
	assert.Equal(t, &parse.Error{Err: ctx.Err()}, err)
}

func TestItems_event(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected parse.Event
	}{
		{
			name: "explicit DTSTART (DATE) and DTEND (DATE)",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTART:20200101
DTEND:20200510
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("DTSTART", "20200101", nil),
					property("DTEND", "20200510", nil),
				},
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2020, time.May, 10, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "implicit DTEND via DURATION",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTART:20200101
DURATION:P2W4D5H2M10S
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("DTSTART", "20200101", nil),
					property("DURATION", "P2W4D5H2M10S", nil),
				},
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
				End: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local).
					AddDate(0, 0, 14).AddDate(0, 0, 4).
					Add(5 * time.Hour).Add(2 * time.Minute).Add(10 * time.Second),
			},
		},
		{
			name: "implicit 1-day duration",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTART:20200101
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("DTSTART", "20200101", nil),
				},
				Start: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local),
				End:   time.Date(2020, time.January, 1, 0, 0, 0, 0, time.Local).AddDate(0, 0, 1),
			},
		},
		{
			name: "implicit until-end-of-day duration",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTART;VALUE=DATE-TIME:20200101T103020
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("DTSTART", "20200101T103020", parse.Parameters{"VALUE": []string{"DATE-TIME"}}),
				},
				Start: time.Date(2020, time.January, 1, 10, 30, 20, 0, time.Local),
				End:   time.Date(2020, time.January, 2, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name: "summary",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
SUMMARY:This is a
  folded summary
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("SUMMARY", "This is a folded summary", nil),
				},
				Summary: "This is a folded summary",
			},
		},
		{
			name: "description",
			input: `BEGIN:VCALENDAR
BEGIN:VEVENT
DESCRIPTION;FMTTYPE=text/plain:A description with a parameter. Also
  folded :)
END:VEVENT
END:VCALENDAR`,
			expected: parse.Event{
				Properties: []parse.Property{
					property("DESCRIPTION", "A description with a parameter. Also folded :)", parse.Parameters{
						"FMTTYPE": []string{"text/plain"},
					}),
				},
				Description: "A description with a parameter. Also folded :)",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cal, err := parse.Items(lex.Text(test.input))
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.expected, cal.Events[0])
		})
	}
}

func item(typ lex.ItemType, val string) lex.Item {
	return lex.Item{
		Type:  typ,
		Value: val,
	}
}

func beginCalendar() lex.Item {
	return item(lex.CalendarBegin, "BEGIN:VCALENDAR")
}

func endCalendar() lex.Item {
	return item(lex.CalendarEnd, "END:VCALENDAR")
}

func beginEvent() lex.Item {
	return item(lex.EventBegin, "BEGIN:VEVENT")
}

func endEvent() lex.Item {
	return item(lex.EventEnd, "END:VEVENT")
}

func property(name, val string, params parse.Parameters) parse.Property {
	if params == nil {
		params = make(parse.Parameters)
	}
	return parse.Property{
		Name:   name,
		Params: params,
		Value:  val,
	}
}

func lexItems(items ...lex.Item) <-chan lex.Item {
	ch := make(chan lex.Item)
	go func() {
		for _, item := range items {
			ch <- item
		}
		close(ch)
	}()
	return ch
}

func loadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
