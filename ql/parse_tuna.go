package ql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MaxHalford/tuna"
)

func newAgg(name, field string) (tuna.Agg, error) {
	agg, ok := map[string]tuna.Agg{
		"MAX":  tuna.NewMax(field),
		"MEAN": tuna.NewMean(field),
	}[strings.ToUpper(name)]
	if !ok {
		return nil, fmt.Errorf("no agg named '%s'", name)
	}
	return agg, nil
}

func parseAgg(lit string) (tuna.Agg, error) {
	// Check the format is <agg>(<field>)
	matched, err := regexp.MatchString("[[:alpha:]]+\\([[:alnum:]]+\\)", lit)
	if !matched || err != nil {
		return nil, fmt.Errorf("couldn't understand '%s', expected a <agg>(<field>)", lit)
	}

	// Parse the Agg name and the field
	parts := strings.Split(lit, "(")
	name, field := parts[0], parts[1][:len(parts[1])-1]

	return newAgg(name, field)
}

func newStream(name, path string) (stream tuna.Stream, err error) {
	switch strings.ToLower(name) {
	case "csv":
		stream, err = tuna.NewCSVStreamFromPath(path)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("no stream named '%s'", name)
	}
	return stream, nil
}

func parseStream(lit string) (tuna.Stream, error) {
	// Check the format is <stream>(<path>)
	matched, err := regexp.MatchString("[[:alpha:]]+\\(.+\\)", lit)
	if !matched || err != nil {
		return nil, fmt.Errorf("couldn't understand '%s', expected something like <stream>(<path>)", lit)
	}

	// Parse the Stream name and the path
	parts := strings.Split(lit, "(")
	name, path := parts[0], parts[1][:len(parts[1])-1]

	return newStream(name, path)
}
