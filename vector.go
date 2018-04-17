package line

// A Vector maps feature positions tZo values. This allows a sparse
// representation of a set of features by only including the non-zero ones.
type Vector map[uint32]float64

// Dot-product between two vectors.
func dot(a, b Vector) (sum float64) {
	// Iterate over the smallest vector
	if len(a) < len(b) {
		for k, v := range a {
			sum += b[k] * v
		}
		return
	}
	for k, v := range b {
		sum += a[k] * v
	}
	return
}
