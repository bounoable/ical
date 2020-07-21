package lex_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/bounoable/ical/internal/testutil"
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
				testutil.BeginCalendar(),
				testutil.Item(lex.Name, "VERSION"),
				testutil.Item(lex.Value, "2.0"),
				testutil.Item(lex.Name, "METHOD"),
				testutil.Item(lex.Value, "REQUEST"),
				testutil.Item(lex.Name, "PRODID"),
				testutil.Item(lex.Value, "Example//Product//ID"),

				testutil.BeginEvent(),
				testutil.Item(lex.Name, "UID"),
				testutil.Item(lex.Value, "111111111111"),
				testutil.Item(lex.Name, "DTSTAMP"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20191010"),
				testutil.Item(lex.Name, "DTSTART"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20200101"),
				testutil.Item(lex.Name, "DTEND"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20200110"),
				testutil.EndEvent(),

				testutil.BeginEvent(),
				testutil.Item(lex.Name, "UID"),
				testutil.Item(lex.Value, "222222222222"),
				testutil.Item(lex.Name, "DTSTAMP"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20191212"),
				testutil.Item(lex.Name, "DTSTART"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20200201"),
				testutil.Item(lex.Name, "DTEND"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20200210"),
				testutil.EndEvent(),

				testutil.EndCalendar(),
				testutil.Item(lex.EOF, ""),
			},
		},
		"ignore invalid linebreaks (LF)": {
			filepath: filepath.Join(wd, "testdata/calendar_lf.ics"),
			expected: []lex.Item{
				testutil.BeginCalendar(),
				testutil.Item(lex.Name, "VERSION"),
				testutil.Item(lex.Value, "2.0"),
				testutil.Item(lex.Name, "METHOD"),
				testutil.Item(lex.Value, "REQUEST"),
				testutil.Item(lex.Name, "PRODID"),
				testutil.Item(lex.Value, "Example//Product//ID"),

				testutil.BeginEvent(),
				testutil.Item(lex.Name, "UID"),
				testutil.Item(lex.Value, "111111111111"),
				testutil.Item(lex.Name, "DTSTAMP"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20191010"),
				testutil.Item(lex.Name, "DTSTART"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20200101"),
				testutil.Item(lex.Name, "DTEND"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20200110"),
				testutil.EndEvent(),

				testutil.BeginEvent(),
				testutil.Item(lex.Name, "UID"),
				testutil.Item(lex.Value, "222222222222"),
				testutil.Item(lex.Name, "DTSTAMP"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20191212"),
				testutil.Item(lex.Name, "DTSTART"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20200201"),
				testutil.Item(lex.Name, "DTEND"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE"),
				testutil.Item(lex.Value, "20200210"),
				testutil.EndEvent(),

				testutil.EndCalendar(),
				testutil.Item(lex.EOF, ""),
			},
		},
		"invalid line breaks (LF) in strict mode": {
			filepath: filepath.Join(wd, "testdata/calendar_lf.ics"),
			opts: []lex.Option{
				lex.StrictLineBreaks,
			},
			expected: []lex.Item{
				testutil.BeginCalendar(),
				testutil.Item(lex.Error, "missing carriage return (CR) at pos 16"),
			},
		},
		"folded (CRLF)": {
			filepath: filepath.Join(wd, "testdata/calendar_folded_crlf.ics"),
			expected: []lex.Item{
				testutil.BeginCalendar(),
				testutil.BeginEvent(),
				testutil.Item(lex.Name, "DESCRIPTION"),
				testutil.Item(lex.Value, "A description that is too long to fit into 75 octets should wrap to the next line. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation."),
				testutil.EndEvent(),
				testutil.EndCalendar(),
				testutil.Item(lex.EOF, ""),
			},
		},
		"folded (LF)": {
			filepath: filepath.Join(wd, "testdata/calendar_folded_lf.ics"),
			expected: []lex.Item{
				testutil.BeginCalendar(),
				testutil.BeginEvent(),
				testutil.Item(lex.Name, "DESCRIPTION"),
				testutil.Item(lex.Value, "A description that is too long to fit into 75 octets should wrap to the next line. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation."),
				testutil.EndEvent(),
				testutil.EndCalendar(),
				testutil.Item(lex.EOF, ""),
			},
		},
		"multiple params": {
			filepath: filepath.Join(wd, "testdata/multiple_params.ics"),
			expected: []lex.Item{
				testutil.BeginCalendar(),
				testutil.BeginEvent(),
				testutil.Item(lex.Name, "ATTACH"),
				testutil.Item(lex.ParamName, "FMTTYPE"),
				testutil.Item(lex.ParamValue, "text/plain"),
				testutil.Item(lex.ParamName, "ENCODING"),
				testutil.Item(lex.ParamValue, "BASE64"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "BINARY"),
				testutil.Item(lex.Value, "VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZy4"),
				testutil.EndEvent(),
				testutil.EndCalendar(),
				testutil.Item(lex.EOF, ""),
			},
		},
		"multiple param values": {
			filepath: filepath.Join(wd, "testdata/multiple_param_values.ics"),
			expected: []lex.Item{
				testutil.BeginCalendar(),
				testutil.BeginEvent(),
				testutil.Item(lex.Name, "X-CUSTOM"),
				testutil.Item(lex.ParamName, "FOO"),
				testutil.Item(lex.ParamValue, "foo"),
				testutil.Item(lex.ParamValue, "bar"),
				testutil.Item(lex.ParamValue, "baz"),
				testutil.Item(lex.ParamName, "BAR"),
				testutil.Item(lex.ParamValue, "baz"),
				testutil.Item(lex.ParamValue, "bar"),
				testutil.Item(lex.ParamValue, "foo"),
				testutil.Item(lex.Value, "foobar"),
				testutil.EndEvent(),
				testutil.EndCalendar(),
				testutil.Item(lex.EOF, ""),
			},
		},
		"with alarm": {
			filepath: filepath.Join(wd, "testdata/with_alarm.ics"),
			expected: []lex.Item{
				testutil.BeginCalendar(),
				testutil.BeginAlarm(),
				testutil.Item(lex.Name, "TRIGGER"),
				testutil.Item(lex.Value, "foo"),
				testutil.Item(lex.Name, "ACTION"),
				testutil.Item(lex.Value, "bar"),
				testutil.EndAlarm(),
				testutil.EndCalendar(),
				testutil.Item(lex.EOF, ""),
			},
		},
		"empty value": {
			filepath: filepath.Join(wd, "testdata/empty_value.ics"),
			expected: []lex.Item{
				testutil.BeginCalendar(),
				testutil.BeginEvent(),
				testutil.Item(lex.Name, "UID"),
				testutil.Item(lex.Value, ""),
				testutil.Item(lex.Name, "DTSTAMP"),
				testutil.Item(lex.ParamName, "VALUE"),
				testutil.Item(lex.ParamValue, "DATE-TIME"),
				testutil.Item(lex.Value, "20191010T000000Z"),
				testutil.Item(lex.Name, "DTSTART"),
				testutil.Item(lex.Value, "20200101"),
				testutil.Item(lex.Name, "DTEND"),
				testutil.Item(lex.Value, "20200110"),
				testutil.EndEvent(),
				testutil.EndCalendar(),
				testutil.Item(lex.EOF, ""),
			},
		},
	}

	for _, test := range tests {
		t.Run(filepath.Base(test.filepath), func(t *testing.T) {
			f, err := os.Open(test.filepath)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			ch := lex.Reader(f, test.opts...)

			var items []lex.Item
			for item := range ch {
				items = append(items, item)
			}

			assert.Equal(t, test.expected, items)
		})
	}
}

func TestLex_context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := lex.File(
		filepath.Join(wd, "testdata/calendar_crlf.ics"),
		lex.Context(ctx),
	)
	if err != nil {
		t.Fatal(err)
	}
	cancel()

	var items []lex.Item
	for item := range ch {
		items = append(items, item)
	}

	assert.Equal(t, lex.Item{
		Type:  lex.Error,
		Value: ctx.Err().Error(),
	}, items[len(items)-1])
}
