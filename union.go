package tuna

// A Union maintains multiple Extractors in parallel.
type Union struct {
	Extractors []Extractor
}

// Update each Extractor given a Row.
func (u *Union) Update(row Row) error {
	for i := range u.Extractors {
		if err := u.Extractors[i].Update(row); err != nil {
			return err
		}
	}
	return nil
}

// Collect concatenates the output of each Extractor's Collect call.
func (u Union) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		row := make(Row)
		for _, ex := range u.Extractors {
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

// Size is the sum of the sizes of each Extractor.
func (u Union) Size() uint {
	var s uint
	for _, ex := range u.Extractors {
		s += ex.Size()
	}
	return s
}

// NewUnion returns a Union with the given Extractors.
func NewUnion(exs ...Extractor) *Union {
	var union = &Union{Extractors: make([]Extractor, len(exs))}
	for i, ex := range exs {
		union.Extractors[i] = ex
	}
	return union
}
