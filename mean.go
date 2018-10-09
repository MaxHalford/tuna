package tuna

import "fmt"

// Mean computes a running average. The result is an approximation but it is
// good enough for most purposes.
type Mean struct {
	Parse  func(Row) (float64, error)
	Prefix string
	n      float64
	mean   float64
}

// Update Mean given a Row.
func (m *Mean) Update(row Row) error {
	var x, err = m.Parse(row)
	if err != nil {
		return err
	}
	m.n++
	m.mean += (x - m.mean) / m.n
	return nil
}

// Collect returns the current value.
func (m Mean) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{fmt.Sprintf("%smean", m.Prefix): float2Str(m.mean)}
		close(c)
	}()
	return c
}

// Size is 1.
func (m Mean) Size() uint { return 1 }

// NewMean returns a Mean that computes the mean of a given field.
func NewMean(field string) *Mean {
	return &Mean{
		Parse:  func(row Row) (float64, error) { return str2Float(row[field]) },
		Prefix: fmt.Sprintf("%s_", field),
	}
}
