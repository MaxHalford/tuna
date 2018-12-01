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
			agg:    NewExtractor("flux", NewSkew()),
			output: "flux_skew\n0\n",
		},
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "5.0"},
			),
			agg:   NewExtractor("fluxx", NewSkew()),
			isErr: true,
		},
	}.Run(t)
}
