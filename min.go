package tuna

import (
	"math"
)

// Min computes the minimal value of a column.
type Min struct {
	min float64
}

// Update Min given a Row.
func (m *Min) Update(x float64) error {
	m.min = math.Min(m.min, x)
	return nil
}

// Collect returns the current value.
func (m Min) Collect() map[string]float64 {
	return map[string]float64{"min": m.min}
}

// NewMin returns a Min that computes the mean of a given field.
func NewMin() *Min {
	return &Min{
		min: math.Inf(1),
	}
}
