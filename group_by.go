package tuna

import (
	"sort"
)

// GroupBy maintains one Agg instance per group.
type GroupBy struct {
	By          string
	NewAgg      func() Agg
	SortResults bool
	groups      map[string]Agg
}

// Update updates the Agg of the Row's group.
func (gb *GroupBy) Update(row Row) error {
	key, ok := row[gb.By]
	if !ok {
		return ErrUnknownField{gb.By}
	}
	if _, ok = gb.groups[key]; !ok {
		gb.groups[key] = gb.NewAgg()
	}
	return gb.groups[key].Update(row)
}

// Collect streams the Collect of each group. The groups are output in the
// lexical order of their keys.
func (gb GroupBy) Collect() <-chan Row {
	// We first list the existing group keys
	keys := make([]string, len(gb.groups))
	var i uint
	for k := range gb.groups {
		keys[i] = k
		i++
	}
	// Sort the group keys so that the output is deterministic
	if gb.SortResults {
		sort.Strings(keys)
	}
	c := make(chan Row)
	go func() {
		for _, key := range keys {
			for r := range gb.groups[key].Collect() {
				// Add the group key to the results
				r[gb.By] = key
				c <- r
			}
		}
		close(c)
	}()
	return c
}

// NewGroupBy returns a GroupBy that maintains a Agg for each
// distinct value of a given variable.
func NewGroupBy(by string, newAgg func() Agg) *GroupBy {
	return &GroupBy{
		By:          by,
		NewAgg:      newAgg,
		SortResults: true,
		groups:      make(map[string]Agg),
	}
}

// SequentialGroupBy maintains one Agg instance. Once a new group key is
// encoutered the Trigger is called. This has many practical use case for large
// but sequential data.
type SequentialGroupBy struct {
	By     string
	NewAgg func() Agg
	Sink   Sink
	key    string
	agg    Agg
}

// Flush writes the results of the Agg and resets it.
func (sgb *SequentialGroupBy) Flush() error {
	if sgb.agg != nil {
		if err := sgb.Sink.Write(sgb.Collect()); err != nil {
			return err
		}
		sgb.agg = sgb.NewAgg()
	} else {
		sgb.agg = sgb.NewAgg()
	}
	return nil
}

// Update updates the Agg of the Row's group.
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
	return sgb.agg.Update(row)
}

// Collect streams the Collect of the current Agg.
func (sgb SequentialGroupBy) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		for r := range sgb.agg.Collect() {
			// Add the group key to the results
			r[sgb.By] = sgb.key
			c <- r
		}
		close(c)
	}()
	return c
}

// NewSequentialGroupBy returns a SequentialGroupBy that maintains an Agg
// for the given variable.
func NewSequentialGroupBy(by string, newAgg func() Agg, sink Sink) *SequentialGroupBy {
	return &SequentialGroupBy{
		By:     by,
		NewAgg: newAgg,
		Sink:   sink,
	}
}
