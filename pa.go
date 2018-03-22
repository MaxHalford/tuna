package main

import "math"

// PAClassifier implements passive-aggressive learning.
type PAClassifier struct {
	c float64
	w Vector
}

// NewPAClassifier instantiates and returns a *PAClassifier.
func NewPAClassifier(c float64) *PAClassifier {
	return &PAClassifier{
		c: c,
		w: make(Vector),
	}
}

// FitPartial of PAClassifier.
func (pa *PAClassifier) FitPartial(x Vector, y float64) {
	// Just a different convention
	if y == 0 {
		y = -1
	}

	var (
		loss  = math.Max(0, 1-y*dot(pa.w, x)) // Hinge loss
		tau   = loss / (dot(x, x) + (1.0 / (2.0 * pa.c)))
		coeff = tau * y
	)

	// Update weights
	for i, xi := range x {
		pa.w[i] += coeff * xi
	}
}

// PredictPartial of PAClassifier.
func (pa PAClassifier) PredictPartial(x Vector) (yHat float64) {
	return sigmoid(dot(pa.w, x))
}
