package tuna

// A Extractor computes a exature in an online manner.
type Extractor interface {
	Update(Row) error
	Collect() <-chan Row
	Size() uint
}
