package tuna

import "testing"

func TestKurtosis(t *testing.T) {
	AggTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "5.0"},
			),
			agg:    NewExtractor("flux", NewKurtosis()),
			output: "flux_kurtosis\n-1.3\n",
		},
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "5.0"},
			),
			agg:   NewExtractor("fluxx", NewKurtosis()),
			isErr: true,
		},
	}.Run(t)
}
