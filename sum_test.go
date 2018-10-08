package tuna

import "testing"

func TestSum(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "4.0"},
				Row{"flux": "-2.0"},
			),
			extractor: NewSum("flux"),
			output:    "flux_sum\n3\n",
		},
	}.Run(t)
}
