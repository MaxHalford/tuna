package tuna

import "fmt"

// Diff runs a Agg on the (x[i+1] - x[i]) version of a stream of values. This
// can be used in conjunction with a GroupBy to compute rolling statistics.
type Diff struct {
	Parse     func(Row) (float64, error)
	Agg       Agg
	FieldName string
	seen      bool
	xi        float64
}

// Update Diff given a Row.
func (d *Diff) Update(row Row) error {
	var x, err = d.Parse(row)
	if err != nil {
		return err
	}
	if !d.seen {
		d.xi = x
		d.seen = true
		return nil
	}
	row[d.FieldName] = float2Str(x - d.xi)
	d.xi = x
	return d.Agg.Update(row)
}

// Collect returns the current value.
func (d Diff) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		for r := range d.Agg.Collect() {
			c <- r
		}
		close(c)
	}()
	return c
}

// Size is the size of the Agg.
func (d Diff) Size() uint { return d.Agg.Size() }

// NewDiff returns a Diff that applies a Agg to the difference of
// a given field.
func NewDiff(field string, newAgg func(s string) Agg) *Diff {
	fn := fmt.Sprintf("%s_diff", field)
	return &Diff{
		Parse:     func(row Row) (float64, error) { return str2Float(row[field]) },
		Agg:       newAgg(fn),
		FieldName: fn,
	}
}
