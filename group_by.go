package tuna

// GroupBy maintains one Extractor instance per group.
type GroupBy struct {
	NewExtractor func() Extractor
	By           string
	groups       map[string]Extractor
}

// Update updates the Extractor of the Row's group.
func (gb *GroupBy) Update(row Row) error {
	key := row[gb.By]
	if _, ok := gb.groups[key]; !ok {
		gb.groups[key] = gb.NewExtractor()
	}
	return gb.groups[key].Update(row)
}

// Collect streams the Collect of each group.
func (gb GroupBy) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		for key, g := range gb.groups {
			for r := range g.Collect() {
				c <- r.Set(gb.By, key)
			}
		}
		close(c)
	}()
	return c
}

// Size is the sum of the sizes of each group.
func (gb GroupBy) Size() uint {
	var s uint
	for _, g := range gb.groups {
		s += g.Size()
	}
	return s
}

// NewGroupBy returns a GroupBy that maintains a Extractor for each
// distinct value of a given variable.
func NewGroupBy(newExtractor func() Extractor, by string) *GroupBy {
	return &GroupBy{
		NewExtractor: newExtractor,
		By:           by,
		groups:       make(map[string]Extractor),
	}
}

// SequentialGroupBy maintains one Extractor instance. Once a new group key is
// encoutered the Trigger is called. This has many practical use case for large
// but sequential data.
type SequentialGroupBy struct {
	NewExtractor func() Extractor
	By           string
	Sink         Sink
	key          string
	extractor    Extractor
}

// Update updates the Extractor of the Row's group.
func (sgb *SequentialGroupBy) Update(row Row) error {
	key := row[sgb.By]
	// Call the Trigger if key has changed
	if sgb.key != key && sgb.extractor != nil {
		if err := sgb.Sink.Write(sgb.Collect()); err != nil {
			return err
		}
		sgb.extractor = sgb.NewExtractor()
	}
	if sgb.extractor == nil {
		sgb.extractor = sgb.NewExtractor()
	}
	sgb.key = key
	return sgb.extractor.Update(row)
}

// Collect streams the Collect of the current Extractor.
func (sgb SequentialGroupBy) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		for r := range sgb.extractor.Collect() {
			c <- r.Set(sgb.By, sgb.key)
		}
		close(c)
	}()
	return c
}

// Size is the size of the current Extractor.
func (sgb SequentialGroupBy) Size() uint {
	return sgb.extractor.Size()
}

// NewSequentialGroupBy returns a SequentialGroupBy that maintains an Extractor
// for the given variable.
func NewSequentialGroupBy(newExtractor func() Extractor, by string, sink Sink) *SequentialGroupBy {
	return &SequentialGroupBy{
		NewExtractor: newExtractor,
		By:           by,
		Sink:         sink,
	}
}
