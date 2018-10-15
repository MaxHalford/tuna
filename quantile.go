package tuna

import (
	"fmt"
	"math/rand"
	"time"
)

const maxHeight = 31

// Quantile computes approximate quantiles using the method presented in
// Greenwald and Khanna 2001.
type Quantile struct {
	Parse     func(Row) (float64, error)
	Prefix    string
	Epsilon   float64
	CutPoints []float64
	sl        *skipList
	n         uint64
}

// Update inserts an item into the quantile sl
func (qu *Quantile) Update(row Row) error {
	var x, err = qu.Parse(row)
	if err != nil {
		return err
	}

	elt := qu.sl.insert(tuple{x, 1, 0})
	qu.n++

	if elt.prev[0] != qu.sl.head && elt.next[0] != nil {
		elt.value.delta = int(2 * qu.Epsilon * float64(qu.n))
	}

	if qu.n%uint64(1.0/float64(2.0*qu.Epsilon)) == 0 {
		qu.compress()
	}

	return nil
}

// Collect returns the estimated quantiles.
func (qu Quantile) Collect() <-chan ErrRow {
	c := make(chan ErrRow)
	go func() {
		r := ErrRow{make(Row), nil}
		for _, cp := range qu.CutPoints {
			r = r.Set(fmt.Sprintf("%sq%s", qu.Prefix, float2Str(cp)), float2Str(qu.Query(cp)))
		}
		c <- r
		close(c)
	}()
	return c
}

// Size is the number of cut points.
func (qu Quantile) Size() uint { return uint(len(qu.CutPoints)) }

// NewQuantile returns a new Quantile with accuracy epsilon (0 <= epsilon <= 1)
func NewQuantile(field string, epsilon float64, cutPoints []float64) *Quantile {
	return &Quantile{
		Parse:     func(row Row) (float64, error) { return str2Float(row[field]) },
		Prefix:    fmt.Sprintf("%s_", field),
		Epsilon:   epsilon,
		CutPoints: cutPoints,
		sl:        newSkipList(),
	}
}

func (qu *Quantile) compress() {

	var missing int

	epsN := int(2 * qu.Epsilon * float64(qu.n))

	for elt := qu.sl.head.next[0]; elt != nil && elt.next[0] != nil; {
		next := elt.next[0]
		t := elt.value
		nt := &next.value

		// value merging
		if t.v == nt.v {
			missing += nt.g
			nt.delta += missing
			nt.g = t.g
			qu.sl.Remove(elt)
		} else if t.g+nt.g+missing+nt.delta < epsN {
			nt.g += t.g + missing
			missing = 0
			qu.sl.Remove(elt)
		} else {
			nt.g += missing
			missing = 0
		}
		elt = next
	}
}

type tuple struct {
	v     float64
	g     int
	delta int
}

// Query returns an epsilon estimate of the element at quantile 'q' (0 <= q <= 1)
func (qu *Quantile) Query(q float64) float64 {

	// convert quantile to rank

	r := int(q*float64(qu.n) + 0.5)

	var rmin int

	epsN := int(qu.Epsilon * float64(qu.n))

	for elt := qu.sl.head.next[0]; elt != nil; elt = elt.next[0] {

		t := elt.value

		rmin += t.g

		n := elt.next[0]

		if n == nil {
			return t.v
		}

		if r+epsN < rmin+n.value.g+n.value.delta {

			if r+epsN < rmin+n.value.g {
				return t.v
			}

			return n.value.v
		}
	}

	panic("not reached")
}

type skipList struct {
	height int
	head   *node
	rnd    *rand.Rand
}

type node struct {
	value tuple
	next  []*node
	prev  []*node
}

func newSkipList() *skipList {
	return &skipList{
		height: 0,
		head:   &node{next: make([]*node, maxHeight)},
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *skipList) insert(t tuple) *node {
	level := 0

	n := s.rnd.Int31()
	for n&1 == 1 {
		level++
		n >>= 1
	}

	if level > s.height {
		s.height++
		level = s.height
	}

	node := &node{
		value: t,
		next:  make([]*node, level+1),
		prev:  make([]*node, level+1),
	}
	curr := s.head
	for i := s.height; i >= 0; i-- {

		for curr.next[i] != nil && t.v >= curr.next[i].value.v {
			curr = curr.next[i]
		}

		if i > level {
			continue
		}

		node.next[i] = curr.next[i]
		if curr.next[i] != nil && curr.next[i].prev[i] != nil {
			curr.next[i].prev[i] = node
		}
		curr.next[i] = node
		node.prev[i] = curr
	}

	return node
}

func (s *skipList) Remove(node *node) {

	// remove n from each level of the skipList

	for i := range node.next {
		prev := node.prev[i]
		next := node.next[i]

		if prev != nil {
			prev.next[i] = next
		}
		if next != nil {
			next.prev[i] = prev
		}
		node.next[i] = nil
		node.prev[i] = nil
	}
}
