package tuna

import (
	"math"
)

// PTP computes the minimal value of a column.
type PTP struct {
	min float64
	max float64
}

// Update PTP given a Row.
func (m *PTP) Update(x float64) error {
	m.min = math.Min(m.min, x)
	m.max = math.Max(m.max, x)
	return nil
}

// Collect returns the current value.
func (m PTP) Collect() map[string]float64 {
	return map[string]float64{"ptp": m.max - m.min}
}

// NewPTP returns a PTP that computes the mean of a given field.
func NewPTP() *PTP {
	return &PTP{
		max: math.Inf(-1),
		min: math.Inf(1),
	}
}
