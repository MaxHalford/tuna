package tuna

import (
	"os"
	"strings"
)

func ExampleRun() {
	// For the sake of example we inline the data, but usually it should be
	// located in a file, database, or some other source
	in := `name,£,bangers
Del Boy,-42,1
Rodney,1001,1
Rodney,1002,2
Del Boy,42,0
Grandad,0,3`

	// Define a Stream
	stream, _ := NewCSVStream(strings.NewReader(in))

	// Define an Agg
	agg := NewGroupBy(
		"name",
		func() Agg {
			return NewUnion(
				NewMean("£"),
				NewSum("bangers"),
			)
		},
	)

	// Define a Sink
	sink := NewCSVSink(os.Stdout)

	// Run
	Run(stream, agg, sink, 0)

	// Output:
	// bangers_sum,name,£_mean
	// 1,Del Boy,0
	// 3,Grandad,0
	// 3,Rodney,1001.5
}
