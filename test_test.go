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
	size      uint
}

// Run runs the Extractor against then Stream and then checks the
// results of the Collect method.
func (tc ExtractorTestCase) Run(t *testing.T) {
	// Go through the Rows and update the Extractor
	Run(tc.stream, tc.extractor, nil, 0)

	// Collect and check the output
	b := &strings.Builder{}
	sink := NewCSVSink(b)
	sink.Write(tc.extractor.Collect())
	output := b.String()
	if output != tc.output {
		t.Errorf("got:\n%swant:\n%s", output, tc.output)
	}

	// Check the size
	size := tc.extractor.Size()
	if size != tc.size {
		t.Errorf("got: %d, want: %d\n", size, tc.size)
	}
}

// ExtractorTestCases is a ExtractorTestCase slice, it's just here for
// convenience.
type ExtractorTestCases []ExtractorTestCase

// Run the test cases.
func (etcs ExtractorTestCases) Run(t *testing.T) {
	for i, tc := range etcs {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) { tc.Run(t) })
	}
}
