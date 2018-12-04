package tuna

import "testing"

func TestVariance(t *testing.T) {
	AggTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "4.0"},
				Row{"flux": "-2.0"},
			),
			agg:    NewExtractor("flux", NewVariance()),
			output: "flux_variance\n6\n",
		},
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "4.0"},
				Row{"flux": "-2.0"},
			),
			agg:   NewExtractor("fluxx", NewVariance()),
			isErr: true,
		},
	}.Run(t)
}
