package tuna

import (
	"math"
)

// Max computes the maximal value of a column.
type Max struct {
	Parse func(Row) (float64, error)
	max   float64
}

// Update Max given a Row.
func (m *Max) Update(row Row) error {
	var x, err = m.Parse(row)
	if err != nil {
		return err
	}
	m.max = math.Max(m.max, x)
	return nil
}

// Collect returns the current value.
func (m Max) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{"max": float64ToString(m.max)}
		close(c)
	}()
	return c
}

// Size is 1.
func (m Max) Size() uint { return 1 }

// NewMax returns a Max that computes the mean of a given field.
func NewMax(field string) *Max {
	return &Max{
		Parse: func(row Row) (float64, error) {
			return stringToFloat64(row[field])
		},
		max: math.Inf(-1),
	}
}
