package tuna

import (
	"fmt"
	"strconv"
)

// An Extractor is the simplest kind of Agg. It feeds the output of Extract to
// a Metric.
type Extractor struct {
	Extract func(row Row) (float64, error)
	Metric  Metric
	Prefix  string
}

// Update parses the Row using Extract and feeds the result to Metric.
func (ex Extractor) Update(row Row) error {
	x, err := ex.Extract(row)
	if err != nil {
		return err
	}
	return ex.Metric.Update(x)
}

// Collect converts the results from the Metric.
func (ex Extractor) Collect() <-chan Row {
	c := make(chan Row)
	go func() {
		row := make(Row)
		for k, v := range ex.Metric.Collect() {
			row[fmt.Sprintf("%s%s", ex.Prefix, k)] = strconv.FormatFloat(v, 'f', -1, 64)
		}
		c <- row
		close(c)
	}()
	return c
}

// NewExtractor returns an Extractor that parses a field as a float64.
func NewExtractor(field string, metrics ...Metric) Extractor {
	return Extractor{
		Extract: func(row Row) (float64, error) { return strconv.ParseFloat(row[field], 64) },
		Prefix:  fmt.Sprintf("%s_", field),
		Metric:  Metrics(metrics),
	}
}
