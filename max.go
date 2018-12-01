package tuna

import (
	"math"
)

// Max computes the maximal value of a column.
type Max struct {
	max float64
}

// Update Max given a Row.
func (m *Max) Update(x float64) error {
	m.max = math.Max(m.max, x)
	return nil
}

// Collect returns the current value.
func (m Max) Collect() map[string]float64 {
	return map[string]float64{"max": m.max}
}

// NewMax returns a Max that computes the mean of a given field.
func NewMax() *Max {
	return &Max{
		max: math.Inf(-1),
	}
}
