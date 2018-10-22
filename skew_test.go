package tuna

import "testing"

func TestSkew(t *testing.T) {
	AggTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "5.0"},
			),
			agg:    NewSkew("flux"),
			output: "flux_skew\n0\n",
			size:   1,
		},
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "5.0"},
			),
			agg:   NewSkew("fluxx"),
			isErr: true,
		},
	}.Run(t)
}
