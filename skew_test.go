package tuna

import "testing"

func TestSkew(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "10.0"},
				Row{"flux": "-4.0"},
			),
			extractor: NewSkew("flux"),
			output:    "flux_skew\n0.43385993540133483\n",
		},
	}.Run(t)
}
