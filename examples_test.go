package tuna

import (
	"os"
	"strings"
)

func Example1() {
	in := `name,£,bangers
"Del Boy",-42,1
Rodney,1001,1
Rodney,1002,2
"Del Boy",42,0
Grandad,0,3`

	// Define a Stream
	stream, _ := NewCSVStream(strings.NewReader(in))

	// Define an Extractor
	extractor := NewGroupBy(
		"name",
		func() Extractor {
			return NewUnion(
				NewMean("£"),
				NewSum("bangers"),
			)
		},
	)

	// Define a Sink
	sink, _ := NewCSVSink(os.Stdout)

	// Run
	Run(stream, extractor, sink, 0)

	// Output:
	// bangers_sum,name,£_mean
	// 1,Del Boy,0
	// 3,Grandad,0
	// 3,Rodney,1001.5
}
