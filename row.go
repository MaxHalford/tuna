package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

// A Row maps a feature name to a raw value.
type Row map[string]string

// A RowIterator returns rows one by one.
type RowIterator interface {
	Read() (row Row, stop bool, err error)
}

// A RowParser parses a row.
type RowParser func(row Row) (ID string, x Vector, y float64)

// A CSVRowReader reads Rows from a CSV file.
type CSVRowReader struct {
	names     []string
	csvReader *csv.Reader
	row       Row
}

// NewCSVRowReader instantiates and returns a CSVRowReader returning rows from
// the given path.
func NewCSVRowReader(path string) (*CSVRowReader, error) {
	// Open the file
	f, err := os.Open(path)
	if err != nil {
		f.Close()
		return nil, err
	}
	// Read the file
	r := csv.NewReader(bufio.NewReader(f))
	// The first row is used for column names
	names, err := r.Read()
	return &CSVRowReader{names, r, make(Row)}, nil
}

// Read of CSVRowReader.
func (reader *CSVRowReader) Read() (row Row, stop bool, err error) {
	r, err := reader.csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, true, err
		}
		return nil, false, err
	}
	for i, name := range reader.names {
		reader.row[name] = r[i]
	}
	return reader.row, false, nil
}
