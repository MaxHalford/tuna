package tuna

import "testing"

func TestDiff(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "4.0"},
			),
			extractor: NewDiff("flux", func(s string) Extractor { return NewMean(s) }),
			output:    "flux_diff_mean\n1.5\n",
			size:      1,
		},
	}.Run(t)
}
