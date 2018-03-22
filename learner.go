package main

import (
	"fmt"
	"time"
)

// A Learner is fed rows one by one and learns progressively.
type Learner interface {
	FitPartial(x Vector, y float64) (yHat float64)
	PredictPartial(x Vector) (yHat float64)
}

func mean(fs []float64) float64 {
	m := 0.0
	for _, f := range fs {
		m += f
	}
	m /= float64(len(fs))
	return m
}

// Fit a stream with a Learner.
func Fit(
	learner Learner,
	stream <-chan Instance,
	valMetric func(y, yHat []float64) float64,
	checkEvery uint64,
) {
	// Store the last checkEvery values in order to obtain an online
	// validation score
	var (
		valYs    = make([]float64, checkEvery)
		valYHats = make([]float64, checkEvery)
		scores   = make([]float64, 0)
		i        int
		start    = time.Now()
	)
	// Loop over the stream of instances
	for instance := range stream {
		// Fit the learning to the instance and obtain the out-of-fold
		// prediction
		yHat := learner.FitPartial(instance.x, instance.y)
		valYHats[i] = yHat
		valYs[i] = instance.y
		i++
		// Monitor progress
		if checkEvery > 0 && instance.t > 0 && (instance.t+1)%checkEvery == 0 {
			// Display the amount of time spent and the number of rows processed
			fmt.Printf("%s -- Processed %d rows", fmtDuration(time.Since(start)), (instance.t + 1))
			if valMetric != nil && len(valYs) > 0 {
				score := valMetric(valYs, valYHats)
				scores = append(scores, score)
				fmt.Printf(" (%f, %f)", score, mean(scores))
			}
			fmt.Print("\n")
			// Reset monitoring variables
			valYs = make([]float64, checkEvery)
			valYHats = make([]float64, checkEvery)
			i = 0
		}
	}
}
