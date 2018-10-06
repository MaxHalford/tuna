package tuna

import (
	"math"
	"strconv"
)

// Skew computes a running skew using an extension of Welford's algorithm.
type Skew struct {
	Parse func(Row) (float64, error)
	n     float64
	mu    float64
	m2    float64
	m3    float64
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
		c <- Row{"skew": strconv.FormatFloat((math.Sqrt(s.n)*s.m3)/math.Pow(s.m2, 1.5), 'f', -1, 64)}
		close(c)
	}()
	return c
}

// Size is 1.
func (s Skew) Size() uint { return 1 }

// NewSkew returns a Skew that computes the mean of a given field.
func NewSkew(field string) *Skew {
	return &Skew{
		Parse: func(row Row) (float64, error) {
			return strconv.ParseFloat(row[field], 64)
		},
	}
}
