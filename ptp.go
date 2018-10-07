package tuna

// PTP computes (max - min). PTP is the acronym for "peak to peak".
type PTP struct {
	Parse func(Row) (float64, error)
	min   *Min
	max   *Max
}

// Update PTP given a Row.
func (ptp *PTP) Update(row Row) error {
	if err := ptp.min.Update(row); err != nil {
		return err
	}
	if err := ptp.max.Update(row); err != nil {
		return err
	}
	return nil
}

// Collect returns the current value.
func (ptp PTP) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{"ptp": float64ToString(ptp.max.max - ptp.min.min)}
		close(c)
	}()
	return c
}

// Size is 1.
func (ptp PTP) Size() uint { return 1 }

// NewPTP returns a PTP that computes the PTP value of a given field.
func NewPTP(field string) *PTP {
	return &PTP{
		Parse: func(row Row) (float64, error) { return stringToFloat64(row[field]) },
		min:   NewMin(field),
		max:   NewMax(field),
	}
}
