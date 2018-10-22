package tuna

import (
	"testing"
)

func TestQuantile(t *testing.T) {
	AggTestCases{
		{
			stream: NewStream(
				Row{"flux": "1.0"},
				Row{"flux": "2.0"},
				Row{"flux": "3.0"},
				Row{"flux": "4.0"},
				Row{"flux": "5.0"},
			),
			agg:    NewQuantile("flux", 0.01, []float64{0.25, 0.5, 0.75}),
			output: "flux_q0.25,flux_q0.5,flux_q0.75\n1,3,4\n",
			size:   3,
		},
	}.Run(t)
}
