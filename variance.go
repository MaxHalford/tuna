package tuna

// Variance computes a running average using Welford's algorithm.
type Variance struct {
	n  float64
	mu float64 // Running mean
	ss float64 // Running sum of squares
}

// Update Variance given a Row.
func (v *Variance) Update(x float64) error {
	v.n++
	// Compute the current mean
	mu := v.mu + (x-v.mu)/v.n
	// Update the sum of squares and the mean
	v.ss += (x - v.mu) * (x - mu)
	v.mu = mu
	return nil
}

// Collect returns the current value.
func (v Variance) Collect() map[string]float64 {
	return map[string]float64{"variance": v.ss / v.n}
}

// NewVariance returns a Variance that computes the mean of a given field.
func NewVariance() *Variance {
	return &Variance{}
}
