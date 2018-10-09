package tuna

import (
	"sort"
)

// GroupBy maintains one Extractor instance per group.
type GroupBy struct {
	By           string
	NewExtractor func() Extractor
	groups       map[string]Extractor
}

// Update updates the Extractor of the Row's group.
func (gb *GroupBy) Update(row Row) error {
	key, ok := row[gb.By]
	if !ok {
		return ErrUnknownField{gb.By}
	}
	if _, ok = gb.groups[key]; !ok {
		gb.groups[key] = gb.NewExtractor()
	}
	return gb.groups[key].Update(row)
}

// Collect streams the Collect of each group. The groups are output in the
// lexical order of their keys.
func (gb GroupBy) Collect() <-chan Row {
	// Sort the group keys so that the output is deterministic
	keys := make([]string, len(gb.groups))
	var i uint
	for k := range gb.groups {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	c := make(chan Row)
	go func() {
		for _, key := range keys {
			for r := range gb.groups[key].Collect() {
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
func NewGroupBy(by string, newExtractor func() Extractor) *GroupBy {
	return &GroupBy{
		By:           by,
		NewExtractor: newExtractor,
		groups:       make(map[string]Extractor),
	}
}

// SequentialGroupBy maintains one Extractor instance. Once a new group key is
// encoutered the Trigger is called. This has many practical use case for large
// but sequential data.
type SequentialGroupBy struct {
	By           string
	NewExtractor func() Extractor
	Sink         Sink
	key          string
	extractor    Extractor
}

// Flush writes the results of the Extractor and resets it.
func (sgb *SequentialGroupBy) Flush() error {
	if sgb.extractor != nil {
		if err := sgb.Sink.Write(sgb.Collect()); err != nil {
			return err
		}
		sgb.extractor = sgb.NewExtractor()
	} else {
		sgb.extractor = sgb.NewExtractor()
	}
	return nil
}

// Update updates the Extractor of the Row's group.
func (sgb *SequentialGroupBy) Update(row Row) error {
	key, ok := row[sgb.By]
	if !ok {
		return ErrUnknownField{sgb.By}
	}
	if sgb.key != key {
		if err := sgb.Flush(); err != nil {
			return err
		}
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
func NewSequentialGroupBy(by string, newExtractor func() Extractor, sink Sink) *SequentialGroupBy {
	return &SequentialGroupBy{
		By:           by,
		NewExtractor: newExtractor,
		Sink:         sink,
	}
}
