package tuna

import "strconv"

// Variance computes a running average using Welford's algorithm.
type Variance struct {
	Parse func(Row) (float64, error)
	n     float64
	mu    float64 // Running mean
	ss    float64 // Running sum of squares
}

// Update Variance given a Row.
func (v *Variance) Update(row Row) error {
	var x, err = v.Parse(row)
	if err != nil {
		return err
	}
	v.n++
	// Compute the current mean
	mu := v.mu + (x-v.mu)/v.n
	// Update the sum of squares and the mean
	v.ss += (x - v.mu) * (x - mu)
	v.mu = mu
	return nil
}

// Collect returns the current value.
func (v Variance) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{"variance": strconv.FormatFloat(v.ss/v.n, 'f', -1, 64)}
		close(c)
	}()
	return c
}

// Size is 1.
func (v Variance) Size() uint { return 1 }

// NewVariance returns a Variance that computes the mean of a given field.
func NewVariance(field string) *Variance {
	return &Variance{
		Parse: func(row Row) (float64, error) { return strconv.ParseFloat(row[field], 64) },
	}
}
