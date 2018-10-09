package tuna

import "testing"

func TestNUnique(t *testing.T) {
	ExtractorTestCases{
		{
			stream: NewStream(
				Row{"ice-cream": "1", "flavor": "chocolate"},
				Row{"ice-cream": "1", "flavor": "coffee"},
				Row{"ice-cream": "1", "flavor": "coffee"},
				Row{"ice-cream": "2", "flavor": "mango"},
				Row{"ice-cream": "2", "flavor": "mango"},
			),
			extractor: NewGroupBy("ice-cream", func() Extractor { return NewNUnique("flavor") }),
			output:    "flavor_n_unique,ice-cream\n2,1\n1,2\n",
			size:      2,
		},
		{
			stream: NewStream(
				Row{"ice-cream": "1", "flavor": "chocolate"},
				Row{"ice-cream": "1", "flavor": "coffee"},
				Row{"ice-cream": "2", "flavor": "mango"},
				Row{"ice-cream": "1", "flavor": "coffee"},
				Row{"ice-cream": "2", "flavor": "mango"},
			),
			extractor: NewGroupBy("ice-cream", func() Extractor { return NewNUnique("flavor") }),
			output:    "flavor_n_unique,ice-cream\n2,1\n1,2\n",
			size:      2,
		},
		{
			stream: NewStream(
				Row{"ice-cream": "1", "flavor": "chocolate"},
				Row{"ice-cream": "1", "flavor": "coffee"},
				Row{"ice-cream": "1", "flavor": "coffee"},
				Row{"ice-cream": "2", "flavor": "mango"},
				Row{"ice-cream": "2", "flavor": "mango"},
			),
			extractor: NewGroupBy("ice-cream", func() Extractor { return NewNUnique("flavorr") }),
			isErr:     true,
		},
	}.Run(t)
}
