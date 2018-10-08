package tuna

import "fmt"

// PTP computes (max - min). PTP is the acronym for "peak to peak".
type PTP struct {
	Parse  func(Row) (float64, error)
	Prefix string
	min    *Min
	max    *Max
}

// Update PTP given a Row.
func (ptp *PTP) Update(row Row) error {
	if err := ptp.min.Update(row); err != nil {
		return err
	}
	return ptp.max.Update(row)
}

// Collect returns the current value.
func (ptp PTP) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		c <- Row{fmt.Sprintf("%s_ptp", ptp.Prefix): float2Str(ptp.max.max - ptp.min.min)}
		close(c)
	}()
	return c
}

// Size is 1.
func (ptp PTP) Size() uint { return 1 }

// NewPTP returns a PTP that computes the PTP value of a given field.
func NewPTP(field string) *PTP {
	return &PTP{
		Parse:  func(row Row) (float64, error) { return str2Float(row[field]) },
		Prefix: field,
		min:    NewMin(field),
		max:    NewMax(field),
	}
}
