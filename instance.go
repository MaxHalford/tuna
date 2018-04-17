package line

import "runtime"

// An Instance in a dataset.
type Instance struct {
	t  uint64
	ID string
	X  Vector
	Y  float64
}

// newInstanceStream returns a channel that sends Instances.
func newInstanceStream(ri RowReader, rp RowParser) <-chan Instance {
	var (
		nCores = runtime.GOMAXPROCS(-1)
		stream = make(chan Instance, nCores*4)
	)
	go func() {
		defer close(stream)
		var t uint64
		for {
			row, stop, _ := ri.Read() // TODO: handle error
			if stop {
				break
			}
			id, x, y := rp(row)
			stream <- Instance{t, id, x, y}
			t++
		}
	}()
	return stream
}
