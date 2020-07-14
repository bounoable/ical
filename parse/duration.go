package parse

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"
)

func parseDuration(raw string) (time.Duration, error) {
	if len(raw) == 0 {
		return 0, nil
	}
	return (&durationParser{value: raw}).parse()
}

type durationParser struct {
	value string
	pos   int
	width int
}

var day = time.Hour * 24
var week = day * 7

func (p *durationParser) parse() (time.Duration, error) {
	r, err := p.next()
	if err != nil {
		return 0, p.unexpectedEnd()
	}

	multiplier := time.Duration(1)

	if r == '-' {
		multiplier = -1
	} else if r != '+' {
		p.backup()
	}

	if r, err = p.next(); r != 'P' {
		return 0, fmt.Errorf("expected 'P' at pos %d; got %s", p.pos+1, string(r))
	}

	var total time.Duration

	for {
		var digits string

		if r, err = p.next(); err != nil {
			return 0, p.unexpectedEnd()
		}

		for unicode.IsDigit(r) {
			digits += string(r)
			if r, err = p.next(); err != nil {
				return 0, p.unexpectedEnd()
			}
		}
		p.backup()

		num, err := strconv.Atoi(digits)
		if err != nil {
			return 0, fmt.Errorf("failed to parse digits in duration: %w", err)
		}

		if r, err = p.next(); err != nil {
			return 0, p.unexpectedEnd()
		}

		var one time.Duration

		switch r {
		case 'W':
			one = week
		case 'D':
			one = day
		case 'H':
			one = time.Hour
		case 'M':
			one = time.Minute
		case 'S':
			one = time.Second
		default:
			return 0, fmt.Errorf("expected one of [W D H M S] at pos %d; got %s", p.pos+1, string(r))
		}

		total += one * time.Duration(num)

		if r, err = p.next(); err != nil {
			return total * multiplier, nil
		}
		p.backup()
	}
}

var errEndOfDuration = errors.New("end of duration")

func (p *durationParser) next() (rune, error) {
	if p.pos >= len(p.value) {
		p.width = 0
		return 0, errEndOfDuration
	}
	r, w := utf8.DecodeRuneInString(p.value[p.pos:])
	p.width = w
	p.pos += w
	return r, nil
}

func (p *durationParser) backup() {
	p.pos -= p.width
}

func (p *durationParser) unexpectedEnd() error {
	return fmt.Errorf("unexpected end of duration at pos %d", p.pos+1)
}
