package line

import (
	"fmt"
	"os"
	"time"
)

// A Model is fed rows one by one and learns progressively.
type Model interface {
	FitPartial(x Vector, y float64) (yPred float64)
	PredictPartial(x Vector) (yPred float64)
}

// Train a Model with a stream of instances.
func Train(
	model Model,
	ri RowReader,
	rp RowParser,
	metric Metric,
	monitor *os.File,
	monitorEvery uint64,
) {
	var (
		stream      = newInstanceStream(ri, rp)
		metricTotal float64
		start       = time.Now()
		i           float64
	)
	// Continue while there are still Instances
	for instance := range stream {
		i++
		// Fit the learning to the instance and obtain the out-of-fold
		// prediction
		yPred := model.FitPartial(instance.X, instance.Y)
		metricTotal += metric.Apply(instance.Y, yPred)
		// Monitor progress
		if monitorEvery > 0 && instance.t > 0 && (instance.t+1)%monitorEvery == 0 {
			duration := time.Since(start)
			monitor.WriteString(fmt.Sprintf(
				"%s -- %d rows -- %d rows/s. -- %f %s\n",
				fmtDuration(time.Since(start)),
				instance.t+1,
				int64(i/duration.Seconds()),
				metricTotal/float64(instance.t+1),
				metric.String(),
			))
			i = 0
		}
	}
}

// Predict the output of a stream of Rows.
func Predict(
	model Model,
	ri RowReader,
	rp RowParser,
	output *os.File,
	monitor *os.File,
	monitorEvery uint64,
) {
	var (
		stream = newInstanceStream(ri, rp)
		start  = time.Now()
	)
	// Continue while there are still Instances
	for instance := range stream {
		yPred := model.PredictPartial(instance.X)
		output.WriteString(fmt.Sprintf("%s,%f\n", instance.ID, yPred))
		// Monitor progress
		if monitorEvery > 0 && instance.t > 0 && (instance.t+1)%monitorEvery == 0 {
			duration := time.Since(start)
			monitor.WriteString(fmt.Sprintf(
				"%s -- %d rows -- %d rows/s.\n",
				fmtDuration(time.Since(start)),
				instance.t+1,
				int64(float64(instance.t+1)/duration.Seconds()),
			))
		}
	}
}
