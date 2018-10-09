package tuna

import "testing"

func TestCount(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"ice-cream": "1", "flavor": "chocolate"},
				Row{"ice-cream": "1", "flavor": "coffee"},
				Row{"ice-cream": "1", "flavor": "vanilla"},
				Row{"ice-cream": "2", "flavor": "mango"},
				Row{"ice-cream": "2", "flavor": "yoghurt"},
			),
			extractor: NewGroupBy("ice-cream", func() Extractor { return NewCount() }),
			output:    "count,ice-cream\n3,1\n2,2\n",
			size:      2,
		},
		{
			stream: NewStream(
				Row{"ice-cream": "1", "flavor": "chocolate"},
				Row{"ice-cream": "1", "flavor": "coffee"},
				Row{"ice-cream": "2", "flavor": "mango"},
				Row{"ice-cream": "1", "flavor": "vanilla"},
				Row{"ice-cream": "2", "flavor": "yoghurt"},
			),
			extractor: NewGroupBy("ice-cream", func() Extractor { return NewCount() }),
			output:    "count,ice-cream\n3,1\n2,2\n",
			size:      2,
		},
	}.Run(t)
}
