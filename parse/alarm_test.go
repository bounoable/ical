package parse_test

import (
	"testing"

	"github.com/bounoable/ical/lex"
	"github.com/bounoable/ical/parse"
	"github.com/stretchr/testify/assert"
)

func TestParse_alarm(t *testing.T) {
	input := `BEGIN:VCALENDAR
BEGIN:VALARM
ACTION:foo
TRIGGER:bar
END:VALARM
END:VCALENDAR`

	items, err := lex.Text(input)
	if err != nil {
		t.Fatal(err)
	}

	cal, err := parse.Items(items)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, cal.Alarms[0], parse.Alarm{
		Properties: []parse.Property{
			property("ACTION", "foo", nil),
			property("TRIGGER", "bar", nil),
		},
		Action:  "foo",
		Trigger: "bar",
	})
}
