package tuna

import (
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Run applies an against a stream. It will display the current progression at
// every multiple of checkpoint.
func Run(stream Stream, extractor Extractor, sink Sink, checkpoint uint) error {
	// Run the exature extractors over the stream
	var (
		n  uint
		t0 = time.Now()
		p  = message.NewPrinter(language.English)
	)
	for row := range stream {
		n++
		// Check there is no error
		if row.err != nil {
			return row.err
		}
		// Update the Extractor
		if err := extractor.Update(row.Row); err != nil {
			return err
		}
		// Display the current progress if a checkpoint is reached
		if checkpoint > 0 && n%checkpoint == 0 {
			t := time.Since(t0)
			p.Printf(
				"\r%s -- %d rows -- %.0f rows/second -- %d values in memory",
				fmtDuration(t),
				n,
				float64(n)/t.Seconds(),
				extractor.Size(),
			)
		}
	}
	// If there was no checkpoint and that there is a sink then the data has to
	// be written
	if checkpoint == 0 && sink != nil {
		if err := sink.Write(extractor.Collect()); err != nil {
			return err
		}
	}
	// If the extractor is a SequentialGroupBy then the last group hasn't been
	// written down yet
	if sgb, ok := extractor.(*SequentialGroupBy); ok {
		if err := sgb.Flush(); err != nil {
			return err
		}
	}
	return nil
}
