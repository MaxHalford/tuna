package tuna

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

// A Stream returns Rows one by one until it's source is depleted.
type Stream chan ErrRow

// StreamRows returns a Stream from a slice of Rows. It is mainly here for
// demonstration and testing purposes.
func StreamRows(rows ...Row) Stream {
	s := make(Stream)
	go func() {
		for _, r := range rows {
			s <- ErrRow{r, nil}
		}
		close(s)
	}()
	return s
}

// StreamCSV returns a Stream from a slice of Rows. It is mainly for
// demonstration and testing purposes.
func StreamCSV(path string) (Stream, error) {

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		file.Close()
		return nil, err
	}

	// Read the file
	reader := csv.NewReader(bufio.NewReader(file))

	// The first row is assumed to contain the column names
	cols, err := reader.Read()
	if err != nil {
		return nil, err
	}

	s := make(Stream)

	go func() {
		for {
			r, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				s <- ErrRow{nil, err}
			}
			tmp := make(Row)
			for i, name := range cols {
				tmp[name] = r[i]
			}
			s <- ErrRow{tmp, nil}
		}
		close(s)
	}()

	return s, nil
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
