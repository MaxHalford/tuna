package tuna

// GroupBy maintains one Extractor instance per group.
type GroupBy struct {
	NewExtractor func() Extractor
	By           string
	groups       map[string]Extractor
}

// Update updates the Extractor of the Row's group.
func (gb GroupBy) Update(row Row) error {
	name := row[gb.By]
	if _, ok := gb.groups[name]; !ok {
		gb.groups[name] = gb.NewExtractor()
	}
	return gb.groups[name].Update(row)
}

// Collect streams the Collect of each group.
func (gb GroupBy) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		for name, g := range gb.groups {
			for r := range g.Collect() {
				c <- r.Set(gb.By, name)
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
