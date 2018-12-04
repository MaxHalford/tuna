package tuna

import "testing"

func TestDiff(t *testing.T) {
	AggTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "4.0"},
			),
			agg:    NewExtractor("flux", NewDiff(NewMean())),
			output: "flux_diff_mean\n1.5\n",
		},
		{
			stream: NewStream(
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "2.0"},
			),
			agg:   NewExtractor("fluxx", NewDiff(NewMean())),
			isErr: true,
		},
	}.Run(t)
}
