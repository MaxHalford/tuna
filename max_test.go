package tuna

import "testing"

func TestMax(t *testing.T) {
	AggTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "4.0"},
				Row{"flux": "2.0"},
			),
			agg:    NewMax("flux"),
			output: "flux_max\n4\n",
			size:   1,
		},
		{
			stream: NewStream(
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "2.0"},
			),
			agg:   NewMax("fluxx"),
			isErr: true,
		},
	}.Run(t)
}
