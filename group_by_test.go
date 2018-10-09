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
			size:      2,
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
			size:      2,
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
			size:      2,
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
			size:   2,
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
					return NewDiff(
						"flux",
						func(s string) Extractor {
							return NewUnion(NewMean(s), NewSum(s))
						},
					)
				},
			),
			output: "flux_diff_mean,flux_diff_sum,key\n1.5,3,a\n0.5,1,b\n",
			size:   4,
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
				"keyy",
				func() Extractor {
					return NewDiff("flux", func(s string) Extractor { return NewMean(s) })
				},
			),
			isErr: true,
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
					return NewDiff(
						"fluxx",
						func(s string) Extractor {
							return NewUnion(NewMean(s), NewSum(s))
						},
					)
				},
			),
			isErr: true,
		},
	}.Run(t)
}

// For SequentialGroupBy we can't use the typical testing boilerplate because
// the output sink has to be provided to NewSequentialGroupBy.
func TestSequentialGroupBy(t *testing.T) {
	var testCases = []struct {
		stream       Stream
		key          string
		newExtractor func() Extractor
		isErr        bool
		output       string
		size         uint
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
			size:         1,
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
			size:         1,
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
			key:          "keyy",
			newExtractor: func() Extractor { return NewMean("flux") },
			isErr:        true,
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
			newExtractor: func() Extractor { return NewMean("fluxx") },
			isErr:        true,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			// Go through the Rows and update the Extractor
			b := &strings.Builder{}
			sgb := NewSequentialGroupBy(tc.key, tc.newExtractor, NewCSVSink(b))
			err := Run(tc.stream, sgb, nil, 0)

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
			output := b.String()
			if output != tc.output {
				t.Errorf("got:\n%swant:\n%s", output, tc.output)
			}

			// Check the size
			size := sgb.Size()
			if size != tc.size {
				t.Errorf("got: %d, want: %d\n", size, tc.size)
			}
		})
	}
}
