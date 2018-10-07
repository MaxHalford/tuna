package tuna

import (
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

// Run applies an against a stream. It will display the current progression at
// every multiple of checkpoint.
func Run(extractor Extractor, stream Stream, writer Writer, checkpoint uint) error {
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
		// Update the Extractor
		if err = extractor.Update(row); err != nil {
			return err
		}
		n++
		if n%checkpoint == 0 {
			// Write the current results
			if writer != nil {
				if err = writer.Write(extractor); err != nil {
					return err
				}
			}
			// Display the current progress
			t := time.Since(t0)
			p.Printf(
				"\r%s -- %d rows -- %.0f rows/second -- %d values",
				fmtDuration(t),
				n,
				float64(n)/t.Seconds(),
				extractor.Size(),
			)
		}
	}
	p.Printf("\nParsed %d rows in %s\n", n, time.Since(t0))
	return nil
}
