package parse

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bounoable/ical/lex"
)

var (
	errEndOfItems = errors.New("end of items")
)

// Calendar ...
type Calendar struct {
	Properties []Property
	ProductID  string
	Version    string
	Calscale   string
	Method     string
	Events     []Event
}

// Event ...
type Event struct {
	Properties  []Property
	UID         string
	Alarms      []Alarm
	Timestamp   time.Time
	Start       time.Time
	End         time.Time
	Summary     string
	Description string
}

// Alarm ...
type Alarm struct {
	Action  string
	Trigger string
}

// Property ...
type Property struct {
	Name   string
	Params Parameters
	Value  string
}

// Parameters ...
type Parameters map[string][]string

// Items ...
func Items(items <-chan lex.Item, opts ...Option) (Calendar, error) {
	p := parser{items: items}
	for _, opt := range opts {
		opt(&p)
	}
	return p.parse()
}

// Slice ...
func Slice(items []lex.Item, opts ...Option) (Calendar, error) {
	ch := make(chan lex.Item)
	go func() {
		defer close(ch)
		for _, item := range items {
			ch <- item
		}
	}()
	return Items(ch, opts...)
}

// Option ...
type Option func(*parser)

// Location ...
func Location(loc *time.Location) Option {
	return func(p *parser) {
		p.loc = loc
	}
}

type parser struct {
	items <-chan lex.Item
	loc   *time.Location

	buf       [2]lex.Item
	start     int
	pos       int
	peekCount int

	cal Calendar
}

func (p *parser) nextItem() (lex.Item, error) {
	item, ok := <-p.items
	if !ok {
		return item, errEndOfItems
	}
	return item, nil
}

func (p *parser) next() (lex.Item, error) {
	if p.peekCount > 0 {
		p.peekCount--
	} else {
		var err error
		if p.buf[0], err = p.nextItem(); err != nil {
			return lex.Item{}, err
		}
	}
	return p.buf[p.peekCount], nil
}

func (p *parser) nextType(typ lex.ItemType) (lex.Item, error) {
	item, err := p.next()
	if err != nil {
		return item, err
	}

	if item.Type != typ {
		return item, p.unexpectedType(item, typ)
	}

	return item, nil
}

func (p *parser) peek() (lex.Item, error) {
	if p.peekCount > 0 {
		return p.buf[p.peekCount-1], nil
	}
	p.peekCount = 1
	var err error
	if p.buf[0], err = p.nextItem(); err != nil {
		return lex.Item{}, err
	}
	return p.buf[0], nil
}

func (p *parser) backup() {
	p.peekCount++
}

func (p *parser) errorf(format string, vals ...interface{}) error {
	return fmt.Errorf(format, vals...)
}

func (p *parser) unexpectedType(item lex.Item, expected lex.ItemType) error {
	return p.errorf("expected item of type %v, but found %s", expected, item)
}

func (p *parser) parse() (Calendar, error) {
	err := p.parseCalendar()
	return p.cal, err
}

func (p *parser) parseCalendar() error {
	item, err := p.next()
	if err != nil {
		return err
	}

	if item.Type != lex.CalendarBegin {
		return p.unexpectedType(item, lex.CalendarBegin)
	}

	cal := Calendar{
		Calscale: "GREGORIAN",
	}

	for {
		item, err = p.next()
		if err != nil {
			return err
		}

		if item.Type == lex.CalendarEnd {
			p.backup()
			break
		}

		if item.Type == lex.EventBegin {
			p.backup()
			evt, err := p.parseEvent()
			if err != nil {
				return err
			}
			cal.Events = append(cal.Events, evt)
		}

		if item.Type == lex.Name {
			p.backup()
			prop, err := p.parseProperty()
			if err != nil {
				return err
			}
			cal.Properties = append(cal.Properties, prop)
		}
	}

	if item.Type != lex.CalendarEnd {
		return p.unexpectedType(item, lex.CalendarEnd)
	}

	for _, prop := range cal.Properties {
		switch prop.Name {
		case "VERSION":
			cal.Version = prop.Value
		case "METHOD":
			cal.Method = prop.Value
		case "PRODID":
			cal.ProductID = prop.Value
		}
	}

	p.cal = cal

	return nil
}

func (p *parser) parseEvent() (Event, error) {
	var evt Event
	item, err := p.nextType(lex.EventBegin)
	if err != nil {
		return evt, err
	}

	for {
		item, err = p.next()
		if err != nil {
			return evt, err
		}

		if item.Type == lex.EventEnd {
			p.backup()
			break
		}

		if item.Type != lex.Name {
			return evt, p.unexpectedType(item, lex.Name)
		}

		p.backup()
		prop, err := p.parseProperty()
		if err != nil {
			return evt, err
		}
		evt.Properties = append(evt.Properties, prop)
	}

	if item, err = p.nextType(lex.EventEnd); err != nil {
		return evt, err
	}

	for _, prop := range evt.Properties {
		switch prop.Name {
		case "UID":
			evt.UID = prop.Value
		case "DTSTART":
			t, err := p.parseTime(prop)
			if err != nil {
				return evt, err
			}
			evt.Start = t
		case "DTEND":
			t, err := p.parseTime(prop)
			if err != nil {
				return evt, err
			}
			evt.End = t
		case "DTSTAMP":
			t, err := p.parseTime(prop)
			if err != nil {
				return evt, err
			}
			evt.Timestamp = t
		}
	}

	return evt, nil
}

func (p *parser) parseProperty() (Property, error) {
	var name string
	params := make(Parameters)

	item, err := p.nextType(lex.Name)
	if err != nil {
		return Property{}, err
	}
	name = item.Value

	if item, err = p.next(); err != nil {
		return Property{}, err
	}

	if item.Type == lex.ParamName {
		p.backup()
		if err = p.parseParams(params); err != nil {
			return Property{}, err
		}
		if item, err = p.nextType(lex.Value); err != nil {
			return Property{}, err
		}
	}

	if item.Type != lex.Value {
		return Property{}, p.unexpectedType(item, lex.Value)
	}

	return Property{
		Name:   name,
		Params: params,
		Value:  item.Value,
	}, nil
}

func (p *parser) parseParams(params Parameters) error {

	for {
		item, err := p.next()
		if err != nil {
			return err
		}

		if item.Type != lex.ParamName {
			p.backup()
			break
		}

		name := item.Value
		var values []string

		for {
			item, err = p.next()
			if err != nil {
				return err
			}

			if item.Type != lex.ParamValue {
				p.backup()
				break
			}

			values = append(values, item.Value)
		}

		params[name] = values
	}

	return nil
}

const (
	layoutDate          = "20060102"
	layoutDateTimeUTC   = "20060102T150405Z"
	layoutDateTimeLocal = "20060102T150405"
)

func (p *parser) parseTime(prop Property) (time.Time, error) {
	var layout string
	loc := time.Local

	if strings.HasSuffix(prop.Value, "Z") {
		layout = layoutDateTimeUTC
		loc = time.UTC
	} else {
		layout = parseLayout(prop.Params)

		if tzRaw, ok := prop.Params["TZID"]; ok {
			for _, raw := range tzRaw {
				if tzloc, err := time.LoadLocation(raw); err == nil {
					loc = tzloc
					break
				}
			}
		}

		if p.loc != nil {
			loc = p.loc
		}
	}

	if layout == layoutDate && len(prop.Value) != len(layoutDate) {
		layout = layoutDateTimeLocal
	}

	return time.ParseInLocation(layout, prop.Value, loc)
}

func parseLayout(params Parameters) string {
	for name, values := range params {
		if name != "VALUE" {
			continue
		}

		for _, val := range values {
			switch val {
			case "DATE":
				return layoutDate
			case "DATE-TIME":
				return layoutDateTimeLocal
			}
		}
	}
	return layoutDate
}
