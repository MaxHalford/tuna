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
			agg:    NewDiff("flux", func(s string) Agg { return NewMean(s) }),
			output: "flux_diff_mean\n1.5\n",
			size:   1,
		},
		{
			stream: NewStream(
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "2.0"},
			),
			agg:   NewDiff("fluxx", func(s string) Agg { return NewMean(s) }),
			isErr: true,
		},
	}.Run(t)
}
