package ql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MaxHalford/tuna"
)

func newExtractor(name, field string) (tuna.Extractor, error) {
	ex, ok := map[string]tuna.Extractor{
		"MAX":  tuna.NewMax(field),
		"MEAN": tuna.NewMean(field),
	}[strings.ToUpper(name)]
	if !ok {
		return nil, fmt.Errorf("no extractor named '%s'", name)
	}
	return ex, nil
}

func parseExtractor(lit string) (tuna.Extractor, error) {
	// Check the format is <extractor>(<field>)
	matched, err := regexp.MatchString("[[:alpha:]]+\\([[:alnum:]]+\\)", lit)
	if !matched || err != nil {
		return nil, fmt.Errorf("couldn't understand '%s', expected a <extractor>(<field>)", lit)
	}

	// Parse the Extractor name and the field
	parts := strings.Split(lit, "(")
	name, field := parts[0], parts[1][:len(parts[1])-1]

	return newExtractor(name, field)
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
