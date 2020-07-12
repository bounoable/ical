package ics

import (
	"io"

	"github.com/bounoable/ical/lex"
	"github.com/bounoable/ical/parse"
)

// Parse ...
func Parse(r io.Reader, opts ...Option) (parse.Calendar, error) {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	items, err := lex.Reader(r, cfg.lexerOptions...)
	if err != nil {
		return parse.Calendar{}, err
	}

	return parse.Items(items, cfg.parserOptions...)
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
