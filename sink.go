package tuna

import (
	"encoding/csv"
	"io"
	"os"
	"sort"
)

// A Sink can persist the output of an Extractor's Collect method.
type Sink interface {
	Write(rows <-chan Row) error
}

// CSVSink persist the output of an Extractor's Collect method to a CSV file.
// The columns are ordered in lexical order.
type CSVSink struct {
	w    *csv.Writer
	cols []string
	tmp  []string
}

// writeRow writes a single Row.
func (cw CSVSink) writeRow(row Row) error {
	for i, c := range cw.cols {
		cw.tmp[i] = row[c]
	}
	return cw.w.Write(cw.tmp)
}

// Write to a CSV located at Path.
func (cw *CSVSink) Write(rows <-chan Row) error {
	defer func() { cw.w.Flush() }()

	if cw.cols == nil {
		// Extract and write the column names and the first row
		cw.cols = make([]string, 0)
		for r := range rows {
			// Extract the columns
			for k := range r {
				cw.cols = append(cw.cols, k)
			}
			// Write the columns
			sort.Strings(cw.cols)
			if err := cw.w.Write(cw.cols); err != nil {
				return err
			}
			// Write the first Row
			cw.tmp = make([]string, len(cw.cols))
			if err := cw.writeRow(r); err != nil {
				return err
			}
			break
		}
	}

	// Write each Row down
	for r := range rows {
		if err := cw.writeRow(r); err != nil {
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
