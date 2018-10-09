package tuna

import "testing"

func TestVariance(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "4.0"},
				Row{"flux": "-2.0"},
			),
			extractor: NewVariance("flux"),
			output:    "flux_variance\n6\n",
			size:      1,
		},
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "4.0"},
				Row{"flux": "-2.0"},
			),
			extractor: NewVariance("fluxx"),
			isErr:     true,
		},
	}.Run(t)
}
