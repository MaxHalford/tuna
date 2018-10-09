package tuna

import "testing"

func TestPTP(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"flux": "3.0"},
				Row{"flux": "4.2"},
				Row{"flux": "-1.0"},
			),
			extractor: NewPTP("flux"),
			output:    "flux_ptp\n5.2\n",
			size:      1,
		},
		{
			stream: NewStream(
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "2.0"},
			),
			extractor: NewPTP("fluxx"),
			isErr:     true,
		},
	}.Run(t)
}
