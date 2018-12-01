package tuna

// Mean computes a running average. The result is an approximation but it is
// good enough for most purposes.
type Mean struct {
	n    float64
	mean float64
}

// Update Mean given a Row.
func (m *Mean) Update(x float64) error {
	m.n++
	m.mean += (x - m.mean) / m.n
	return nil
}

// Collect returns the current value.
func (m Mean) Collect() map[string]float64 {
	return map[string]float64{"mean": m.mean}
}

// NewMean returns a Mean.
func NewMean() *Mean {
	return &Mean{}
}
