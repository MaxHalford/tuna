package tuna

import (
	"fmt"
	"math"
)

// Skew computes a running skew using an extension of Welford's algorithm.
type Skew struct {
	Parse  func(Row) (float64, error)
	Prefix string
	n      float64
	mu     float64
	m2     float64
	m3     float64
}

// Update Skew given a Row.
func (s *Skew) Update(row Row) error {
	var x, err = s.Parse(row)
	if err != nil {
		return err
	}
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
func (s Skew) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{fmt.Sprintf("%sskew", s.Prefix): float2Str((math.Sqrt(s.n) * s.m3) / math.Pow(s.m2, 1.5))}
		close(c)
	}()
	return c
}

// Size is 1.
func (s Skew) Size() uint { return 1 }

// NewSkew returns a Skew that computes the mean of a given field.
func NewSkew(field string) *Skew {
	return &Skew{
		Parse:  func(row Row) (float64, error) { return str2Float(row[field]) },
		Prefix: fmt.Sprintf("%s_", field),
	}
}
