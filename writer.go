package tuna

import (
	"encoding/csv"
	"os"
)

// A Writer can persist the output of an Extractor's Collect method.
type Writer interface {
	Write(ex Extractor) error
}

// CSVWriter persist the output of an Extractor's Collect method to a CSV file.
type CSVWriter struct {
	Path   string
	cols   []string
	tmp    []string
	writer *csv.Writer
}

// Write to a CSV located at Path.
func (cw *CSVWriter) Write(ex Extractor) error {
	if cw.writer == nil {

		// Initialize the file and the associated Writer
		file, err := os.Create(cw.Path)
		if err != nil {
			return err
		}
		cw.writer = csv.NewWriter(file)
		defer cw.writer.Flush()

		// Extract and persist the column names and the first Row
		res := ex.Collect()
		cw.cols = make([]string, 0)
		cw.tmp = make([]string, 0)
		for r := range res {
			for k, v := range r {
				cw.cols = append(cw.cols, k)
				cw.tmp = append(cw.tmp, v)
			}
			if err := cw.writer.Write(cw.cols); err != nil {
				return err
			}
			if err := cw.writer.Write(cw.tmp); err != nil {
				return err
			}
			break
		}
	}

	// Write each Row down
	for r := range ex.Collect() {
		for i, c := range cw.cols {
			cw.tmp[i] = r[c]
		}
		if err := cw.writer.Write(cw.tmp); err != nil {
			return err
		}
	}

	return nil
}

// NewCSVWriter returns a CSVWriter which persists results to the given path.
func NewCSVWriter(path string) *CSVWriter {
	return &CSVWriter{Path: path}
}
