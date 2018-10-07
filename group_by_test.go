package tuna

import "testing"

func TestGroupBy(t *testing.T) {
	ExtractorTestCases{
		{
			stream: StreamRows(
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "3.0"},
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
			),
			extractor: NewGroupBy("key", func() Extractor { return NewMean("flux") }),
			output:    "key,mean\na,2\nb,-2\n",
		},
		{
			stream: StreamRows(
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "3.0"},
			),
			extractor: NewGroupBy("key", func() Extractor { return NewMean("flux") }),
			output:    "key,mean\na,2\nb,-2\n",
		},
		{
			stream: StreamRows(
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "a", "flux": "3.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "-3.0"},
				Row{"key": "a", "flux": "2.0"},
			),
			extractor: NewGroupBy("key", func() Extractor { return NewMean("flux") }),
			output:    "key,mean\na,2\nb,-2\n",
		},
		{
			stream: StreamRows(
				Row{"key": "a", "flux": "1.0"},
				Row{"key": "a", "flux": "2.0"},
				Row{"key": "a", "flux": "4.0"},
				Row{"key": "b", "flux": "-1.0"},
				Row{"key": "b", "flux": "-2.0"},
				Row{"key": "b", "flux": "0.0"},
			),
			extractor: NewGroupBy(
				"key",
				func() Extractor {
					return NewDiff("flux", func(s string) Extractor { return NewMean(s) })
				},
			),
			output: "diff_mean,key\n1.5,a\n0.5,b\n",
		},
	}.Run(t)
}
