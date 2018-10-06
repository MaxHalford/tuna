package tuna

import (
	"encoding/csv"
	"os"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// A Extractor computes a exature in an online manner.
type Extractor interface {
	Update(Row) error
	Collect() <-chan Row
	Size() uint
}

// Run applies a set of Extractors against a stream. It will display the
// current progression at every multiple of notify.
func Run(exs map[string]Extractor, stream Stream, notify uint) error {
	// Run the exature extractors over the stream
	var (
		n  uint
		t0 = time.Now()
		p  = message.NewPrinter(language.English)
	)
	for {
		row, stop, err := stream.Next()
		// Check for stopage
		if stop {
			break
		}
		// Check for error
		if err != nil {
			return err
		}
		// Update each Extractor
		for _, ex := range exs {
			if err = ex.Update(row); err != nil {
				return err
			}
		}
		// Display the current progress
		n++
		if n%notify == 0 {
			var size uint
			for _, ex := range exs {
				size += ex.Size()
			}
			p.Printf(
				"\r%.0f rows/second -- %d rows -- %d values",
				float64(n)/time.Since(t0).Seconds(),
				n,
				size,
			)
		}
	}
	p.Printf("\nParsed %d rows in %s\n", n, time.Since(t0))
	return nil
}

// ToCSV saves the results of a Extractor to a CSV file.
func ToCSV(ex Extractor, name string) error {
	// Create the output file and an associated writer
	file, err := os.Create(name)
	defer file.Close()
	if err != nil {
		return err
	}
	w := csv.NewWriter(file)
	defer w.Flush()

	// Extract the column names
	res := ex.Collect()
	var cols = make([]string, 0)
	for r := range res {
		for c := range r {
			cols = append(cols, c)
		}
		w.Write(cols)
		break
	}

	// Write down the rows one by one
	var tmp = make([]string, len(cols))
	for r := range res {
		for i, c := range cols {
			tmp[i] = r[c]
		}
		if err := w.Write(tmp); err != nil {
			return err
		}
	}

	return nil
}
