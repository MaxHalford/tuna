package tuna

// Sum computes a running sum.
type Sum struct {
	sum float64
}

// Update Sum given a Row.
func (s *Sum) Update(x float64) error {
	s.sum += x
	return nil
}

// Collect returns the current value.
func (s Sum) Collect() map[string]float64 {
	return map[string]float64{"sum": s.sum}
}

// NewSum returns a Sum that computes the mean of a given field.
func NewSum() *Sum {
	return &Sum{}
}
