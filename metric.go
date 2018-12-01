package tuna

// A Metric takes float64s in and spits out a map of float64s.
type Metric interface {
	Update(x float64) error
	Collect() map[string]float64
}

// Metrics is also a Metric.
type Metrics []Metric

// Update calls Update on each element.
func (ms Metrics) Update(x float64) error {
	for i := range ms {
		if err := ms[i].Update(x); err != nil {
			return err
		}
	}
	return nil
}

// Collect merges the results from each element.
func (ms Metrics) Collect() map[string]float64 {
	var r = make(map[string]float64)
	for _, m := range ms {
		for k, v := range m.Collect() {
			r[k] = v
		}
	}
	return r
}
