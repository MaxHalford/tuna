package tuna

// A Agg computes a feature in an online manner.
type Agg interface {
	Update(Row) error
	Collect() <-chan Row
	Size() uint
}
