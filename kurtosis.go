package tuna

import (
	"fmt"
)

// Kurtosis computes a running kurtosis using an extension of Welford's
// algorithm.
type Kurtosis struct {
	Parse  func(Row) (float64, error)
	Prefix string
	n      float64
	mu     float64
	m2     float64
	m3     float64
	m4     float64
}

// Update Kurtosis given a Row.
func (k *Kurtosis) Update(row Row) error {
	var x, err = k.Parse(row)
	if err != nil {
		return err
	}
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
func (k Kurtosis) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{fmt.Sprintf("%skurtosis", k.Prefix): float2Str((k.n*k.m4)/(k.m2*k.m2) - 3)}
		close(c)
	}()
	return c
}

// Size is 1.
func (k Kurtosis) Size() uint { return 1 }

// NewKurtosis returns a Kurtosis that computes the mean of a given field.
func NewKurtosis(field string) *Kurtosis {
	return &Kurtosis{
		Parse:  func(row Row) (float64, error) { return str2Float(row[field]) },
		Prefix: fmt.Sprintf("%s_", field),
	}
}
