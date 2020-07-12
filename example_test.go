package ical_test

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bounoable/ical"
	"github.com/bounoable/ical/lex"
	"github.com/bounoable/ical/parse"
)

func ExampleParse() {
	f, err := os.Open("/path/to/calendar.ics")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	cal, err := ical.Parse(f)
	if err != nil {
		panic(err)
	}

	// work with the events
	for _, evt := range cal.Events {
		fmt.Println(evt)
	}
}

func ExampleParse_withContext() {
	f, err := os.Open("/path/to/calendar.ics")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ctx := context.TODO()
	cal, err := ical.Parse(f, ical.Context(ctx))
	if err != nil {
		panic(err)
	}

	// work with the events
	for _, evt := range cal.Events {
		fmt.Println(evt)
	}
}

func ExampleParse_withLexerOptions() {
	f, err := os.Open("/path/to/calendar.ics")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	lexContext := context.TODO()

	cal, err := ical.Parse(f, ical.LexWith(
		lex.Context(lexContext),
		lex.StrictLineBreaks,
	))

	if err != nil {
		panic(err)
	}

	// work with the events
	for _, evt := range cal.Events {
		fmt.Println(evt)
	}
}

func ExampleParse_withParserOptions() {
	f, err := os.Open("/path/to/calendar.ics")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	parseContext := context.TODO()

	cal, err := ical.Parse(f, ical.ParseWith(
		parse.Context(parseContext),
		parse.Location(time.Local),
	))

	if err != nil {
		panic(err)
	}

	// work with the events
	for _, evt := range cal.Events {
		fmt.Println(evt)
	}
}

func ExampleParse_explicitTimeLocation() {
	f, err := os.Open("/path/to/calendar.ics")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	cal, err := ical.Parse(f, ical.ParseWith(
		parse.Location(loc),
	))

	if err != nil {
		panic(err)
	}

	// work with the events
	for _, evt := range cal.Events {
		fmt.Println(evt)
	}
}
