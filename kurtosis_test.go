package tuna

import (
	"fmt"
	"testing"
)

// An ExtractorTestCase is a generic way to test the output of a Extractor.
type ExtractorTestCase struct {
	fe     Extractor
	stream Stream
	result string
}

// Run runs the Extractor against then Stream and then checks the
// results of the Collect method.
func (tc ExtractorTestCase) Run(t *testing.T) {
	for {
		row, stop, err := tc.stream.Next()
		// Check for stopage
		if stop {
			break
		}
		// Check for error
		if err != nil {
			t.Error(err)
		}
		// Update the Extractor
		if err = tc.fe.Update(row); err != nil {
			t.Error(err)
		}
	}
	// Check the output
	for r := range tc.fe.Collect() {
		fmt.Println(r)
	}
}

func TestKurtosis(t *testing.T) {
	var testCases = []ExtractorTestCase{
		{
			fe: NewKurtosis("flux"),
			stream: &RowStream{[]Row{
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "10.0"},
				Row{"flux": "-4.0"},
			}},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TC %d", i), func(t *testing.T) { tc.Run(t) })
	}
}
