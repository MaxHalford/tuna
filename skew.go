package tuna

import (
	"math"
)

// Skew computes a running skew using an extension of Welford's algorithm.
type Skew struct {
	n  float64
	mu float64
	m2 float64
	m3 float64
}

// Update Skew given a Row.
func (s *Skew) Update(x float64) error {
	s.n++
	delta := x - s.mu
	deltaN := delta / s.n
	term1 := delta * deltaN * (s.n - 1)
	s.mu += deltaN
	s.m3 += term1*deltaN*(s.n-2) - 3*deltaN*s.m2
	s.m2 += term1
	return nil
}

// Collect returns the current value.
func (s Skew) Collect() map[string]float64 {
	return map[string]float64{"skew": (math.Sqrt(s.n) * s.m3) / math.Pow(s.m2, 1.5)}
}

// NewSkew returns a Skew that computes the mean of a given field.
func NewSkew() *Skew {
	return &Skew{}
}
