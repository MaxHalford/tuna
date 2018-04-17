package line

import "math"

// A Metric can compute an online metric.
type Metric interface {
	Apply(yTrue, yPred float64) float64
	String() string
}

// LogLoss returns -yTrue * log(yPred) - (1-yTrue) * log(1-yPred) where yTrue
// is 0 or 1 and yPred is a probability.
type LogLoss struct{}

// Apply of LogLoss.
func (ll LogLoss) Apply(yTrue, yPred float64) float64 {
	yPred = clip(yPred, 0.00001)
	return -yTrue*math.Log(yPred) - (1-yTrue)*math.Log(1-yPred)
}

// String of LogLoss.
func (ll LogLoss) String() string {
	return "log loss"
}

// Accuracy returns 1 if abs(yTrue - yPred) < 0.5, if not it returns 0.
type Accuracy struct{}

// Apply of Accuracy.
func (acc Accuracy) Apply(yTrue, yPred float64) float64 {
	if math.Abs(yTrue-yPred) < 0.5 {
		return 1
	}
	return 0
}

// String of Accuracy
func (acc Accuracy) String() string {
	return "accuracy"
}
