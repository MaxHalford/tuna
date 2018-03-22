package main

import "math"

// FTRLProximalClassifier implements the Follow The Regularized Leader
// algorithm described in
// http://www.eecs.tufts.edu/~dsculley/papers/ad-click-prediction.pdf.
type FTRLProximalClassifier struct {
	alpha float64
	beta  float64
	l1    float64
	l2    float64

	n map[string]float64 // Squared sum of past gradients
	z map[string]float64 // Weights
	w map[string]float64 // Lazy weights
}

// NewFTRLProximalClassifier instantiates and returns a *FTRLProximalClassifier.
func NewFTRLProximalClassifier(alpha, beta, l1, l2 float64) *FTRLProximalClassifier {
	return &FTRLProximalClassifier{
		alpha: alpha,
		beta:  beta,
		l1:    l1,
		l2:    l2,
		n:     make(map[string]float64),
		z:     make(map[string]float64),
		w:     make(map[string]float64),
	}
}

// FitPartial of FTRLProximalClassifier.
func (ftrl *FTRLProximalClassifier) FitPartial(x Vector, y float64) (yHat float64) {
	yHat = ftrl.PredictPartial(x)
	g := yHat - y // Gradient
	g2 := g * g
	// Update z and n
	for i := range x {
		sigma := (math.Sqrt(ftrl.n[i]+g2) - math.Sqrt(ftrl.n[i])) / ftrl.alpha
		ftrl.z[i] += g - sigma*ftrl.w[i]
		ftrl.n[i] += g2
	}
	return
}

// PredictPartial of FTRLProximalClassifier.
func (ftrl *FTRLProximalClassifier) PredictPartial(x Vector) (yHat float64) {
	var (
		wTx float64 // Inner product
		w   = make(map[string]float64)
	)
	for i := range x {
		// Get sign of weight z[i]
		z, ok := ftrl.z[i]
		if !ok {
			ftrl.z[i] = 0.0
			ftrl.n[i] = 0.0
		}
		sign := 1.0
		if z < 0 {
			sign = -1.0
		}
		// Build w[i] on the fly using z[i] and n[i]
		if sign*z <= ftrl.l1 {
			w[i] = 0 // l1 regularization
		} else {
			w[i] = (sign*ftrl.l1 - z) / ((ftrl.beta+math.Sqrt(ftrl.n[i]))/ftrl.alpha + ftrl.l2)
			wTx += w[i]
		}
	}
	ftrl.w = w
	return sigmoid(wTx)
}
