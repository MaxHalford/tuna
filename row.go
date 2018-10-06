package tuna

import "fmt"

// A Row maps a column name to a raw value.
type Row map[string]string

// Set set the field of a Row and then returns the Row.
func (r Row) Set(k string, v string) Row {
	r[k] = v
	return r
}

// Suffix adds a suffix to each field of a Row and then returns the Row.
func (r Row) Suffix(suffix string, sep string) Row {
	var nr = make(map[string]string)
	for k, v := range r {
		nr[fmt.Sprintf("%s%s%s", k, sep, suffix)] = v
	}
	return nr
}
