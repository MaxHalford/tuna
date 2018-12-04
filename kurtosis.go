package tuna

// Kurtosis computes a running kurtosis using an extension of Welford's
// algorithm.
type Kurtosis struct {
	n  float64
	mu float64
	m2 float64
	m3 float64
	m4 float64
}

// Update Kurtosis given a Row.
func (k *Kurtosis) Update(x float64) error {
	k.n++
	delta := x - k.mu
	deltaN := delta / k.n
	deltaN2 := deltaN * deltaN
	term1 := delta * deltaN * (k.n - 1)
	k.mu += deltaN
	k.m4 += term1*deltaN2*(k.n*k.n-3*k.n+3) + 6*deltaN2*k.m2 - 4*deltaN*k.m3
	k.m3 += term1*deltaN*(k.n-2) - 3*deltaN*k.m2
	k.m2 += term1
	return nil
}

// Collect returns the current value.
func (k Kurtosis) Collect() map[string]float64 {
	return map[string]float64{"kurtosis": (k.n*k.m4)/(k.m2*k.m2) - 3}
}

// NewKurtosis returns a Kurtosis.
func NewKurtosis() *Kurtosis {
	return &Kurtosis{}
}
