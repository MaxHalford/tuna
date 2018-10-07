package tuna

import (
	"encoding/csv"
	"os"
	"sort"
)

// A Sink can persist the output of an Extractor's Collect method.
type Sink interface {
	Write(rows <-chan Row) error
}

// CSVSink persist the output of an Extractor's Collect method to a CSV file.
type CSVSink struct {
	Path   string
	cols   []string
	tmp    []string
	writer *csv.Writer
}

// Write to a CSV located at Path.
func (cw *CSVSink) Write(rows <-chan Row) error {
	if cw.writer == nil {

		// Initialize the file and the associated Sink
		file, err := os.Create(cw.Path)
		if err != nil {
			return err
		}
		cw.writer = csv.NewWriter(file)
		defer cw.writer.Flush()

		// Extract and persist the column names and the first Row
		cw.cols = make([]string, 0)
		for r := range rows {
			for k := range r {
				cw.cols = append(cw.cols, k)
			}
			sort.Strings(cw.cols)
			if err := cw.writer.Write(cw.cols); err != nil {
				return err
			}
			break
		}
		cw.tmp = make([]string, len(cw.cols))
	}

	// Write each Row down
	for r := range rows {
		for i, c := range cw.cols {
			cw.tmp[i] = r[c]
		}
		if err := cw.writer.Write(cw.tmp); err != nil {
			return err
		}
	}

	return nil
}

// NewCSVSink returns a CSVSink which persists results to the given path.
func NewCSVSink(path string) *CSVSink {
	return &CSVSink{Path: path}
}
