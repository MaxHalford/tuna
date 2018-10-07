package tuna

import (
	"encoding/csv"
	"io"
	"os"
)

// A Sink can persist the output of an Extractor's Collect method.
type Sink interface {
	Write(rows <-chan Row) error
}

// CSVSink persist the output of an Extractor's Collect method to a CSV file.
type CSVSink struct {
	w    *csv.Writer
	cols []string
	tmp  []string
}

// Write to a CSV located at Path.
func (cw *CSVSink) Write(rows <-chan Row) error {
	defer func() { cw.w.Flush() }()
	if cw.cols == nil {
		// Extract and write the column names and the first row
		cw.cols = make([]string, 0)
		cw.tmp = make([]string, 0)
		for r := range rows {
			for k, v := range r {
				cw.cols = append(cw.cols, k)
				cw.tmp = append(cw.tmp, v)
			}
			if err := cw.w.Write(cw.cols); err != nil {
				return err
			}
			if err := cw.w.Write(cw.tmp); err != nil {
				return err
			}
			break
		}
	}

	// Write each Row down
	for r := range rows {
		for i, c := range cw.cols {
			cw.tmp[i] = r[c]
		}
		if err := cw.w.Write(cw.tmp); err != nil {
			return err
		}
	}

	return nil
}

// NewCSVSink returns a CSVSink which persists results to the given file.
func NewCSVSink(writer io.Writer) (*CSVSink, error) {
	return &CSVSink{w: csv.NewWriter(writer)}, nil
}

// NewCSVSinkFromPath returns a CSVSink which persists results to the given
// path.
func NewCSVSinkFromPath(path string) (*CSVSink, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return NewCSVSink(file)
}
