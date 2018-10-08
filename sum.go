package tuna

import "fmt"

// Sum computes a running sum.
type Sum struct {
	Parse  func(Row) (float64, error)
	Prefix string
	sum    float64
}

// Update Sum given a Row.
func (s *Sum) Update(row Row) error {
	var x, err = s.Parse(row)
	if err != nil {
		return err
	}
	s.sum += x
	return nil
}

// Collect returns the current value.
func (s Sum) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{fmt.Sprintf("%s_sum", s.Prefix): float2Str(s.sum)}
		close(c)
	}()
	return c
}

// Size is 1.
func (s Sum) Size() uint { return 1 }

// NewSum returns a Sum that computes the mean of a given field.
func NewSum(field string) *Sum {
	return &Sum{
		Parse:  func(row Row) (float64, error) { return str2Float(row[field]) },
		Prefix: field,
	}
}
