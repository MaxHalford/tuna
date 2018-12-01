package tuna

import "fmt"

// Diff runs an Agg on the (x[i+1] - x[i]) version of a stream of values. This
// can be used in conjunction with a GroupBy to compute rolling statistics.
type Diff struct {
	Metric Metric
	seen   bool
	xi     float64
}

// Update Diff given a Row.
func (d *Diff) Update(x float64) error {
	if !d.seen {
		d.xi = x
		d.seen = true
		return nil
	}
	defer func() { d.xi = x }()
	return d.Metric.Update(x - d.xi)
}

// Collect returns the current value.
func (d Diff) Collect() map[string]float64 {
	var r = make(map[string]float64)
	for k, v := range d.Metric.Collect() {
		r[fmt.Sprintf("diff_%s", k)] = v
	}
	return r
}

// NewDiff returns a Diff.
func NewDiff(m Metric) *Diff {
	return &Diff{Metric: m}
}
