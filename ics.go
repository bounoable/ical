package ics

import (
	"io"

	"github.com/bounoable/ical/lex"
	"github.com/bounoable/ical/parse"
)

// Calendar ...
type Calendar parse.Calendar

// Parse ...
func Parse(r io.Reader, opts ...Option) (Calendar, error) {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	items, err := lex.Reader(r, cfg.lexerOptions...)
	if err != nil {
		return Calendar{}, err
	}

	cal, err := parse.Items(items, cfg.parserOptions...)
	if err != nil {
		return Calendar{}, err
	}

	return Calendar(cal), nil
}

// Option ...
type Option func(*config)

// LexWith ...
func LexWith(opts ...lex.Option) Option {
	return func(cfg *config) {
		cfg.lexerOptions = append(cfg.lexerOptions, opts...)
	}
}

// ParseWith ...
func ParseWith(opts ...parse.Option) Option {
	return func(cfg *config) {
		cfg.parserOptions = append(cfg.parserOptions, opts...)
	}
}

type config struct {
	lexerOptions  []lex.Option
	parserOptions []parse.Option
}
