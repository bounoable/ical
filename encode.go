package ical

import (
	"io"

	"github.com/bounoable/ical/encode"
	"github.com/bounoable/ical/parse"
)

// Encode writes the .ics file for cal into w.
func Encode(cal Calendar, w io.Writer) error {
	return encode.Calendar(parse.Calendar(cal), w)
}
