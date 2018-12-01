package tuna

import "testing"

func TestMin(t *testing.T) {
	AggTestCases{
		{
			stream: NewStream(
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "2.0"},
			),
			agg:    NewExtractor("flux", NewMin()),
			output: "flux_min\n2\n",
		},
		{
			stream: NewStream(
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "2.0"},
			),
			agg:   NewExtractor("fluxx", NewMin()),
			isErr: true,
		},
	}.Run(t)
}
