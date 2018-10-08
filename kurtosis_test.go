package tuna

import "testing"

func TestKurtosis(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "10.0"},
				Row{"flux": "-4.0"},
			),
			extractor: NewKurtosis("flux"),
			output:    "flux_kurtosis\n-0.9761404848253483\n",
		},
	}.Run(t)
}
