package lex_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bounoable/ical/lex"
	"github.com/stretchr/testify/assert"
)

var wd, _ = os.Getwd()

func TestReader(t *testing.T) {
	tests := map[string]struct {
		filepath string
		opts     []lex.Option
		expected []lex.Item
	}{
		"valid ical": {
			filepath: filepath.Join(wd, "testdata/calendar_crlf.ics"),
			expected: []lex.Item{
				beginCalendar(),
				item(lex.Name, "VERSION"),
				item(lex.Value, "2.0"),
				item(lex.Name, "METHOD"),
				item(lex.Value, "REQUEST"),
				item(lex.Name, "PRODID"),
				item(lex.Value, "Example//Product//ID"),

				beginEvent(),
				item(lex.Name, "UID"),
				item(lex.Value, "111111111111"),
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20191010"),
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
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20191212"),
				item(lex.Name, "DTSTART"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20200201"),
				item(lex.Name, "DTEND"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20200210"),
				endEvent(),

				endCalendar(),
				item(lex.EOF, ""),
			},
		},
		"ignore invalid linebreaks (LF)": {
			filepath: filepath.Join(wd, "testdata/calendar_lf.ics"),
			expected: []lex.Item{
				beginCalendar(),
				item(lex.Name, "VERSION"),
				item(lex.Value, "2.0"),
				item(lex.Name, "METHOD"),
				item(lex.Value, "REQUEST"),
				item(lex.Name, "PRODID"),
				item(lex.Value, "Example//Product//ID"),

				beginEvent(),
				item(lex.Name, "UID"),
				item(lex.Value, "111111111111"),
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20191010"),
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
				item(lex.Name, "DTSTAMP"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20191212"),
				item(lex.Name, "DTSTART"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20200201"),
				item(lex.Name, "DTEND"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "DATE"),
				item(lex.Value, "20200210"),
				endEvent(),

				endCalendar(),
				item(lex.EOF, ""),
			},
		},
		"invalid line breaks (LF) in strict mode": {
			filepath: filepath.Join(wd, "testdata/calendar_lf.ics"),
			opts: []lex.Option{
				lex.StrictLineBreaks,
			},
			expected: []lex.Item{
				beginCalendar(),
				item(lex.Error, "missing carriage return (CR) at pos 16"),
			},
		},
		"folded (CRLF)": {
			filepath: filepath.Join(wd, "testdata/calendar_folded_crlf.ics"),
			expected: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "DESCRIPTION"),
				item(lex.Value, "A description that is too long to fit into 75 octets should wrap to the next line. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation."),
				endEvent(),
				endCalendar(),
				item(lex.EOF, ""),
			},
		},
		"folded (LF)": {
			filepath: filepath.Join(wd, "testdata/calendar_folded_lf.ics"),
			expected: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "DESCRIPTION"),
				item(lex.Value, "A description that is too long to fit into 75 octets should wrap to the next line. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation."),
				endEvent(),
				endCalendar(),
				item(lex.EOF, ""),
			},
		},
		"multiple params": {
			filepath: filepath.Join(wd, "testdata/multiple_params.ics"),
			expected: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "ATTACH"),
				item(lex.ParamName, "FMTTYPE"),
				item(lex.ParamValue, "text/plain"),
				item(lex.ParamName, "ENCODING"),
				item(lex.ParamValue, "BASE64"),
				item(lex.ParamName, "VALUE"),
				item(lex.ParamValue, "BINARY"),
				item(lex.Value, "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZy4"),
				endEvent(),
				endCalendar(),
				item(lex.EOF, ""),
			},
		},
		"multiple param values": {
			filepath: filepath.Join(wd, "testdata/multiple_param_values.ics"),
			expected: []lex.Item{
				beginCalendar(),
				beginEvent(),
				item(lex.Name, "X-CUSTOM"),
				item(lex.ParamName, "FOO"),
				item(lex.ParamValue, "foo"),
				item(lex.ParamValue, "bar"),
				item(lex.ParamValue, "baz"),
				item(lex.ParamName, "BAR"),
				item(lex.ParamValue, "baz"),
				item(lex.ParamValue, "bar"),
				item(lex.ParamValue, "foo"),
				item(lex.Value, "foobar"),
				endEvent(),
				endCalendar(),
				item(lex.EOF, ""),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.filepath, func(t *testing.T) {
			f, err := os.Open(test.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			ch, err := lex.Reader(f, test.opts...)
			assert.Nil(t, err)

			var items []lex.Item

			for item := range ch {
				items = append(items, item)
			}

			assert.Equal(t, test.expected, items)
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
