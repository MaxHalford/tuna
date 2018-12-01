package tuna

import "testing"

func TestSum(t *testing.T) {
	AggTestCases{
		{
			stream: ZipStreams(
				NewStream(
					Row{"flux": "1.0"},
					Row{"flux": "4.0"},
					Row{"flux": "-2.0"},
				),
				NewStream(
					Row{"flux": "1.0"},
					Row{"flux": "2.0"},
					Row{"flux": "3.0"},
				),
			),
			agg:    NewSum("flux"),
			output: "flux_sum\n9\n",
			size:   1,
		},
		{
			stream: ZipStreams(
				NewStream(
					Row{"fluxx": "1.0"},
					Row{"fluxx": "4.0"},
					Row{"fluxx": "-2.0"},
				),
				NewStream(
					Row{"flux": "1.0"},
					Row{"flux": "2.0"},
					Row{"flux": "3.0"},
				),
			),
			agg:   NewSum("fluxx"),
			isErr: true,
		},
	}.Run(t)
}
