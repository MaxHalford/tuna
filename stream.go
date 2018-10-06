package tuna

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

// A Stream returns Rows one by one until it's source is depleted.
type Stream interface {
	Next() (row Row, stop bool, err error)
}

// A RowStream directly streams Rows. It's mostly here is for testing purposes.
type RowStream struct {
	rows []Row
}

// Next returns the next Row.
func (rs *RowStream) Next() (row Row, stop bool, err error) {
	if len(rs.rows) == 0 {
		return nil, true, nil
	}
	r := rs.rows[0]
	rs.rows = append(rs.rows[:0], rs.rows[1:]...)
	return r, false, nil
}

// A CSVStream reads Rows from a CSV file.
type CSVStream struct {
	names  []string
	reader *csv.Reader
	row    Row
}

// Next returns the next Row.
func (cs *CSVStream) Next() (row Row, stop bool, err error) {
	r, err := cs.reader.Read()
	if err != nil {
		if err == io.EOF {
			return nil, true, err
		}
		return nil, false, err
	}
	for i, name := range cs.names {
		cs.row[name] = r[i]
	}
	return cs.row, false, nil
}

// NewCSVStream returns a CSVStream that streams from the given file.
func NewCSVStream(path string) (*CSVStream, error) {
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
	return &CSVStream{names, r, make(Row)}, nil
}

// A StreamZip iterates over multiple streams one by one. This is rather
// practical for iterating over a single file which has been split into
// multiple smaller files.
type StreamZip struct {
	Streams []Stream
	i       int
}

// Next returns the next row of the next Stream.
func (sz StreamZip) Next() (row Row, stop bool, err error) {
	// Stop if there are no more streams to go through
	if sz.i == len(sz.Streams) {
		return nil, true, nil
	}
	// Read the next line in the current stream
	row, stop, err = sz.Streams[sz.i].Next()
	if stop {
		sz.i++
		return sz.Next()
	}
	return
}
