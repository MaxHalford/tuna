package tuna

// An Agg takes Rows in and spits Rows out.
type Agg interface {
	Update(Row) error
	Collect() <-chan Row
}

// Aggs is also an Agg.
type Aggs []Agg

// Update calls Update on each Agg.
func (aggs Aggs) Update(row Row) error {
	for _, agg := range aggs {
		if err := agg.Update(row); err != nil {
			return err
		}
	}
	return nil
}

// Collect merges and returns the outputs of each Agg's Collect method.
func (aggs Aggs) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		row := make(Row)
		for _, agg := range aggs {
			for r := range agg.Collect() {
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
