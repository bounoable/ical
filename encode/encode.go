package encode

import (
	"io"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/bounoable/ical/parse"
)

// Calendar writes the .ics file for cal into w.
func Calendar(cal parse.Calendar, w io.Writer) error {
	return NewEncoder(w).Encode(cal)
}

// NewEncoder returns a new Encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// Encoder writes .ics files.
type Encoder struct {
	w io.Writer
}

// Encode writes cal as a .ics file to the writer.
func (enc *Encoder) Encode(cal parse.Calendar) error {
	var err error

	if err = enc.string("BEGIN:VCALENDAR"); err != nil {
		return err
	}

	for _, prop := range cal.Properties {
		if err = enc.property(prop); err != nil {
			return err
		}
	}

	for _, evt := range cal.Events {
		if err = enc.event(evt); err != nil {
			return err
		}
	}

	if err = enc.string("\r\nEND:VCALENDAR"); err != nil {
		return err
	}

	return nil
}

func (enc *Encoder) write(p []byte) (int, error) {
	return enc.w.Write(p)
}

func (enc *Encoder) string(s string) error {
	_, err := enc.w.Write([]byte(s))
	return err
}

func (enc *Encoder) property(prop parse.Property) error {
	type parameter struct {
		name   string
		values []string
	}

	var err error
	var linebuilder strings.Builder
	linebuilder.WriteString(prop.Name)

	params := make([]parameter, 0, len(prop.Params))
	for name, vals := range prop.Params {
		params = append(params, parameter{
			name:   name,
			values: vals,
		})
	}

	sort.Slice(params, func(a, b int) bool { return params[a].name < params[b].name })

	for _, param := range params {
		if _, err = linebuilder.WriteString(";" + param.name); err != nil {
			return err
		}
		valstr := strings.Join(param.values, ",")
		if _, err = linebuilder.WriteString("=" + valstr); err != nil {
			return err
		}
	}

	if _, err = linebuilder.WriteString(":" + prop.Value); err != nil {
		return err
	}

	line := linebuilder.String()

	var splits []string

	var l, r int
	for l, r = 0, 75; r < len(line); l, r = r, r+75 {
		for !utf8.RuneStart(line[r]) {
			r--
		}
		splits = append(splits, line[l:r])
	}
	splits = append(splits, line[l:])

	line = "\r\n" + strings.Join(splits, "\r\n ")

	return enc.string(line)
}

func (enc *Encoder) event(evt parse.Event) error {
	var err error
	if err = enc.string("\r\nBEGIN:VEVENT"); err != nil {
		return err
	}

	for _, prop := range evt.Properties {
		if err = enc.property(prop); err != nil {
			return err
		}
	}

	return enc.string("\r\nEND:VEVENT")
}
