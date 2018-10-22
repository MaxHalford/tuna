package tuna

import (
	"strconv"
)

// Count computes a running kurtosis using an extension of Welford's algorithm.
type Count struct {
	n uint64
}

// Update Count given a Row.
func (c *Count) Update(row Row) error {
	c.n++
	return nil
}

// Collect returns the current value.
func (c Count) Collect() <-chan Row {
	ch := make(chan Row)
	go func() {
		ch <- Row{"count": strconv.FormatUint(c.n, 10)}
		close(ch)
	}()
	return ch
}

// Size is 1.
func (c Count) Size() uint { return 1 }

// NewCount returns a Count that computes the mean of a given field.
func NewCount() *Count {
	return &Count{0}
}
