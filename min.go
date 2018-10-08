package tuna

import (
	"fmt"
	"math"
)

// Min computes the minimal value of a column.
type Min struct {
	Parse  func(Row) (float64, error)
	Prefix string
	min    float64
}

// Update Min given a Row.
func (m *Min) Update(row Row) error {
	var x, err = m.Parse(row)
	if err != nil {
		return err
	}
	m.min = math.Min(m.min, x)
	return nil
}

// Collect returns the current value.
func (m Min) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{fmt.Sprintf("%s_min", m.Prefix): float2Str(m.min)}
		close(c)
	}()
	return c
}

// Size is 1.
func (m Min) Size() uint { return 1 }

// NewMin returns a Min that computes the mean of a given field.
func NewMin(field string) *Min {
	return &Min{
		Parse:  func(row Row) (float64, error) { return str2Float(row[field]) },
		Prefix: field,
		min:    math.Inf(1),
	}
}
