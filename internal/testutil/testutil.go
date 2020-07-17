package testutil

import (
	"time"

	"github.com/bounoable/ical/lex"
	"github.com/bounoable/ical/parse"
)

// Item creates a lexer item.
func Item(typ lex.ItemType, val string) lex.Item {
	return lex.Item{
		Type:  typ,
		Value: val,
	}
}

// BeginCalendar creates a lex.CalendarBegin item.
func BeginCalendar() lex.Item {
	return Item(lex.CalendarBegin, "BEGIN:VCALENDAR")
}

// EndCalendar creates a lex.CalendarEnd item.
func EndCalendar() lex.Item {
	return Item(lex.CalendarEnd, "END:VCALENDAR")
}

// BeginEvent creates a lex.EventBegin item.
func BeginEvent() lex.Item {
	return Item(lex.EventBegin, "BEGIN:VEVENT")
}

// EndEvent creates a lex.EventEnd item.
func EndEvent() lex.Item {
	return Item(lex.EventEnd, "END:VEVENT")
}

// BeginAlarm creates a lex.AlarmBegin item.
func BeginAlarm() lex.Item {
	return Item(lex.AlarmBegin, "BEGIN:VALARM")
}

// EndAlarm creates a lex.AlarmEnd item.
func EndAlarm() lex.Item {
	return Item(lex.AlarmEnd, "END:VALARM")
}

// Property creates a parse.Property.
func Property(name, val string, params parse.Parameters) parse.Property {
	if params == nil {
		params = make(parse.Parameters)
	}
	return parse.Property{
		Name:   name,
		Params: params,
		Value:  val,
	}
}

// LexItems lexes a slice of items.
func LexItems(items ...lex.Item) <-chan lex.Item {
	ch := make(chan lex.Item)
	go func() {
		for _, item := range items {
			ch <- item
		}
		close(ch)
	}()
	return ch
}

// LoadLocation loads a time.Location and panics if it fails.
func LoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}
