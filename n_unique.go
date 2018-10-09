package tuna

import (
	"fmt"
	"strconv"
)

// NUnique computes a running sum.
type NUnique struct {
	Parse  func(Row) (string, error)
	Prefix string
	seen   map[string]bool
}

// Update NUnique given a Row.
func (nu *NUnique) Update(row Row) error {
	var x, err = nu.Parse(row)
	if err != nil {
		return err
	}
	nu.seen[x] = true
	return nil
}

// Collect returns the current value.
func (nu NUnique) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{fmt.Sprintf("%sn_unique", nu.Prefix): strconv.Itoa(len(nu.seen))}
		close(c)
	}()
	return c
}

// Size is 1.
func (nu NUnique) Size() uint { return 1 }

// NewNUnique returns a NUnique that computes the mean of a given field.
func NewNUnique(field string) *NUnique {
	return &NUnique{
		Parse:  func(row Row) (string, error) { return row.Get(field) },
		Prefix: fmt.Sprintf("%s_", field),
		seen:   make(map[string]bool),
	}
}
