package main

import (
	"fmt"
	"math"
	"time"
)

// Dot-product between two vectors.
func dot(a, b Vector) (sum float64) {
	// Iterate over the smallest vector
	if len(a) < len(b) {
		for k, v := range a {
			sum += b[k] * v
		}
	} else {
		for k, v := range b {
			sum += a[k] * v
		}
	}
	return sum
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}
