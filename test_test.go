package tuna

import (
	"fmt"
	"strings"
	"testing"
)

// An ExtractorTestCase is a generic way to test the output of a Extractor.
type ExtractorTestCase struct {
	stream    Stream
	extractor Extractor
	output    string
}

// Run runs the Extractor against then Stream and then checks the
// results of the Collect method.
func (tc ExtractorTestCase) Run(t *testing.T) {
	// Go through the Rows and update the Extractor
	for row := range tc.stream {
		if row.err != nil {
			t.Error(row.err)
		}
		if err := tc.extractor.Update(row.Row); err != nil {
			t.Error(err)
		}
	}

	// Collect the output
	b := &strings.Builder{}
	sink, err := NewCSVSink(b)
	if err != nil {
		t.Error(err)
	}
	sink.Write(tc.extractor.Collect())

	// Check the output
	output := b.String()
	if output != tc.output {
		t.Errorf("got:\n%swant:\n%s", output, tc.output)
	}
}

// ExtractorTestCases is a ExtractorTestCase slice, it's just here for
// convenience.
type ExtractorTestCases []ExtractorTestCase

// Run the test cases.
func (etcs ExtractorTestCases) Run(t *testing.T) {
	for i, tc := range etcs {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) { tc.Run(t) })
	}
}
