package tuna

import (
	"fmt"
	"strings"
	"testing"
)

func TestGroupBy(t *testing.T) {
	AggTestCases{
		{
			stream: NewStream(
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "3.0"},
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
			),
			agg:    NewGroupBy("key", func() Agg { return NewExtractor("flux", NewMean()) }),
			output: "flux_mean,key\n2,a\n-2,b\n",
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
			agg:    NewGroupBy("key", func() Agg { return NewExtractor("flux", NewMean()) }),
			output: "flux_mean,key\n2,a\n-2,b\n",
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
			agg:    NewGroupBy("key", func() Agg { return NewExtractor("flux", NewMean()) }),
			output: "flux_mean,key\n2,a\n-2,b\n",
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
			agg: NewGroupBy(
				"key",
				func() Agg { return NewExtractor("flux", NewDiff(NewMean())) },
			),
			output: "flux_diff_mean,key\n1.5,a\n0.5,b\n",
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
			agg: NewGroupBy(
				"key",
				func() Agg { return NewExtractor("flux", NewDiff(Metrics{NewMean(), NewSum()})) },
			),
			output: "flux_diff_mean,flux_diff_sum,key\n1.5,3,a\n0.5,1,b\n",
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
			agg: NewGroupBy(
				"keyy",
				func() Agg { return NewExtractor("flux", NewDiff(NewMean())) },
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
			agg: NewGroupBy(
				"key",
				func() Agg { return NewExtractor("fluxx", NewDiff(Metrics{NewMean(), NewSum()})) },
			),
			isErr: true,
		},
	}.Run(t)
}

// For SequentialGroupBy we can't use the typical testing boilerplate because
// the output sink has to be provided to NewSequentialGroupBy.
func TestSequentialGroupBy(t *testing.T) {
	var testCases = []struct {
		stream Stream
		key    string
		newAgg func() Agg
		isErr  bool
		output string
		size   uint
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
			key:    "key",
			newAgg: func() Agg { return NewExtractor("flux", NewMean()) },
			output: "flux_mean,key\n2,a\n-2,b\n",
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
			key:    "key",
			newAgg: func() Agg { return NewExtractor("flux", NewMean()) },
			output: "flux_mean,key\n-2,b\n2,a\n",
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
			key:    "keyy",
			newAgg: func() Agg { return NewExtractor("flux", NewMean()) },
			isErr:  true,
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
			key:    "key",
			newAgg: func() Agg { return NewExtractor("fluxx", NewMean()) },
			isErr:  true,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			// Go through the Rows and update the Agg
			b := &strings.Builder{}
			sgb := NewSequentialGroupBy(tc.key, tc.newAgg, NewCSVSink(b))
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
		})
	}
}
