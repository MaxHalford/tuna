package main

import "sort"

// A Vector maps feature names to values. This allows a sparse representation
// of a set of features.
type Vector map[string]float64

// sortedKeys returns the keys of a Vector sorted in lexicographical order.
func (vec Vector) sortedKeys() []string {
	var (
		keys = make([]string, len(vec))
		i    int
	)
	for k := range vec {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}
