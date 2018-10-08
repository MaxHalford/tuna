package tuna

import (
	"fmt"
	"strings"
	"testing"
)

func TestGroupBy(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "3.0"},
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
			),
			extractor: NewGroupBy("key", func() Extractor { return NewMean("flux") }),
			output:    "flux_mean,key\n2,a\n-2,b\n",
		},
		{
			stream: NewStream(
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "3.0"},
			),
			extractor: NewGroupBy("key", func() Extractor { return NewMean("flux") }),
			output:    "flux_mean,key\n2,a\n-2,b\n",
		},
		{
			stream: NewStream(
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "a", "flux": "3.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
				Row{"key": "a", "flux": "2.0"},
			),
			extractor: NewGroupBy("key", func() Extractor { return NewMean("flux") }),
			output:    "flux_mean,key\n2,a\n-2,b\n",
		},
		{
			stream: NewStream(
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "4.0"},
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "0.0"},
			),
			extractor: NewGroupBy(
				"key",
				func() Extractor {
					return NewDiff("flux", func(s string) Extractor { return NewMean(s) })
				},
			),
			output: "flux_diff_mean,key\n1.5,a\n0.5,b\n",
		},
	}.Run(t)
}

// For SequentialGroupBy we can't use the typical testing boilerplate because
// the output sink has to be provided to NewSequentialGroupBy.
func TestSequentialGroupBy(t *testing.T) {
	var testCases = []struct {
		stream       Stream
		b            *strings.Builder
		key          string
		newExtractor func() Extractor
		output       string
	}{
		{
			stream: NewStream(
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "3.0"},
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
			),
			key:          "key",
			newExtractor: func() Extractor { return NewMean("flux") },
			output:       "flux_mean,key\n2,a\n-2,b\n",
		},
		{
			stream: NewStream(
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "3.0"},
			),
			key:          "key",
			newExtractor: func() Extractor { return NewMean("flux") },
			output:       "flux_mean,key\n-2,b\n2,a\n",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			// Go through the Rows and update the Extractor
			b := &strings.Builder{}
			sgb := NewSequentialGroupBy(tc.key, tc.newExtractor, NewCSVSink(b))
			Run(tc.stream, sgb, nil, 0)

			// Collect the output
			output := b.String()

			// Check the output
			if output != tc.output {
				t.Errorf("got:\n%swant:\n%s", output, tc.output)
			}
		})
	}
}
