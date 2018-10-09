package tuna

import "testing"

func TestKurtosis(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "5.0"},
			),
			extractor: NewKurtosis("flux"),
			output:    "flux_kurtosis\n-1.3\n",
			size:      1,
		},
	}.Run(t)
}
