package tuna

// A Union maintains multiple Aggs in parallel.
type Union struct {
	Aggs []Agg
}

// Update each Agg given a Row.
func (u *Union) Update(row Row) error {
	for i := range u.Aggs {
		if err := u.Aggs[i].Update(row); err != nil {
			return err
		}
	}
	return nil
}

// Collect concatenates the output of each Agg's Collect call.
func (u Union) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		row := make(Row)
		for _, ex := range u.Aggs {
			for r := range ex.Collect() {
				for k, v := range r {
					row[k] = v
				}
			}
		}
		c <- row
		close(c)
	}()
	return c
}

// Size is the sum of the sizes of each Agg.
func (u Union) Size() uint {
	var s uint
	for _, ex := range u.Aggs {
		s += ex.Size()
	}
	return s
}

// NewUnion returns a Union with the given Aggs.
func NewUnion(exs ...Agg) *Union {
	var union = &Union{Aggs: make([]Agg, len(exs))}
	for i, ex := range exs {
		union.Aggs[i] = ex
	}
	return union
}
