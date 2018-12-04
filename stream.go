package tuna

import (
	"encoding/csv"
	"io"
	"os"
)

// A Stream returns Rows one by one until it's source is depleted.
type Stream chan ErrRow

// NewStream returns a Stream from a slice of Rows. It is mainly here for
// demonstration and testing purposes.
func NewStream(rows ...Row) Stream {
	s := make(Stream)
	go func() {
		for _, r := range rows {
			s <- ErrRow{r, nil}
		}
		close(s)
	}()
	return s
}

// NewFuncStream returns a Stream that calls function n times and returns the
// resulting Rows.
func NewFuncStream(f func() Row, n uint) Stream {
	s := make(Stream)
	go func() {
		for i := uint(0); i < n; i++ {
			s <- ErrRow{f(), nil}
		}
		close(s)
	}()
	return s
}

// NewCSVStream returns a Stream from an io.Reader that reads strings that are
// assumed to CSV-parsable.
func NewCSVStream(reader io.Reader) (Stream, error) {
	csvr := csv.NewReader(reader)

	// The first row is assumed to contain the column names
	cols, err := csvr.Read()
	if err != nil {
		return nil, err
	}

	s := make(Stream)

	go func() {
		for {
			r, err := csvr.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				s <- ErrRow{nil, err}
			}
			row := make(Row)
			for i, name := range cols {
				row[name] = r[i]
			}
			s <- ErrRow{row, nil}
		}
		close(s)
	}()

	return s, nil
}

// NewCSVStreamFromPath returns a Stream from a CSV file.
func NewCSVStreamFromPath(path string) (Stream, error) {
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return nil, err
	}
	return NewCSVStream(file)
}

// ZipStreams returns a Stream that iterates over multiple streams one by one.
// This is quite convinient for going through a dataset which has been split
// into multiple parts.
func ZipStreams(ss ...Stream) Stream {
	zs := make(Stream)
	go func() {
		for _, s := range ss {
			for r := range s {
				zs <- r
			}
		}
		close(zs)
	}()
	return zs
}
