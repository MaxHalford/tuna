package tuna

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func BenchmarkRun(b *testing.B) {
	var (
		letters = strings.Split("abcdefghijklmnopqrstuvwxyz", "")
		randRow = func() Row {
			return Row{
				"letter": letters[rand.Intn(len(letters))],
				"m1":     fmt.Sprintf("%f", rand.Float64()),
				"m2":     fmt.Sprintf("%f", rand.Float64()),
				"m3":     fmt.Sprintf("%f", rand.Float64()),
			}
		}
		stream = NewFuncStream(randRow, 10000)
		agg    = NewGroupBy(
			"letter",
			func() Agg {
				return Aggs{
					NewExtractor("m0", NewMean()),
					NewExtractor("m1", NewMean()),
					NewExtractor("m2", NewMean()),
				}
			},
		)
	)
	for n := 0; n < b.N; n++ {
		Run(stream, agg, nil, 0)
	}
}
