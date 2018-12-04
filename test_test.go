package tuna

import (
	"fmt"
	"strings"
	"testing"
)

// An AggTestCase is a generic way to test the output of a Agg.
type AggTestCase struct {
	stream Stream
	agg    Agg
	isErr  bool
	output string
}

// Run runs the Metric against then Stream and then checks the
// results of the Collect method.
func (tc AggTestCase) Run(t *testing.T) {
	// Go through the Rows and update the Agg
	err := Run(tc.stream, tc.agg, nil, 0)

	// Check the error
	if err == nil {
		if tc.isErr == true {
			t.Error("expected an error, got nil")
			return
		}
	} else {
		if tc.isErr == false {
			t.Errorf("expected no error, got %v", err)
			return
		}
		return
	}

	// Collect and check the output
	b := &strings.Builder{}
	sink := NewCSVSink(b)
	sink.Write(tc.agg.Collect())
	output := b.String()
	if output != tc.output {
		t.Errorf("got:\n%swant:\n%s", output, tc.output)
	}
}

// AggTestCases is a AggTestCase slice, it's just here for
// convenience.
type AggTestCases []AggTestCase

// Run the test cases.
func (etcs AggTestCases) Run(t *testing.T) {
	for i, tc := range etcs {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) { tc.Run(t) })
	}
}
