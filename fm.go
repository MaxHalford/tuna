package main

import "math"

// FMClassifier implements the Factorization Machine algorithm described in
// Rendle's paper (https://www.csie.ntu.edu.tw/~b97053/paper/Rendle2010FM.pdf).
type FMClassifier struct {
	k   uint32  // Number of latent vectors
	eta float64 // Learning rate
	w0  float64 // Global bias
	c0  float64
	w   Vector // Weight of each variable
	c   Vector
	ww  map[uint32]Vector // Weight of variable i in factor f
}

// NewFMClassifier instantiates and returns a *FMClassifier.
func NewFMClassifier(k uint32, eta float64) *FMClassifier {
	// Initialize weights
	var ww = make(map[uint32]Vector)
	for f := uint32(0); f < k; f++ {
		ww[f] = make(Vector)
	}
	return &FMClassifier{
		k:   k,
		eta: eta,
		w0:  0.0,
		c0:  0.0,
		w:   make(Vector),
		c:   make(Vector),
		ww:  ww,
	}
}

// FitPartial of FMClassifier with Stochastic Gradient Descent. The gradients
// are indicated in equation 4 of Rendle's paper.
func (fm *FMClassifier) FitPartial(x Vector, y float64) (yHat float64) {

	yHat = fm.PredictPartial(x)

	v := make(map[uint32]float64)
	for i, xi := range x {
		for f := uint32(0); f < fm.k; f++ {
			v[f] += fm.ww[f][i] * xi
		}
	}

	g := yHat - y
	g2 := g * g

	fm.w0 -= fm.eta / (math.Sqrt(fm.c0) + 1) * g

	for i, xi := range x {
		d := fm.eta / (math.Sqrt(fm.c[i]) + 1) * g * xi
		fm.w[i] -= d
		for f := uint32(0); f < fm.k; f++ {
			fm.ww[f][i] -= d * (v[f] - fm.ww[f][i]*xi)
		}
		fm.c[i] += g2
	}
	fm.c0 += g2

	return
}

// PredictPartial of FTRLProximalClassifier. This is equation 1 of Rendle's
// paper.
func (fm *FMClassifier) PredictPartial(x Vector) (yHat float64) {

	yHat += fm.w0

	v := make(map[uint32]float64)
	vv := make(map[uint32]float64)
	for i, xi := range x {
		yHat += fm.w[i] * x[i]
		for f := uint32(0); f < fm.k; f++ {
			v[f] += fm.ww[f][i] * xi
			vv[f] += (math.Pow(fm.ww[f][i], 2)) * (math.Pow(xi, 2))
		}
	}

	for f := uint32(0); f < fm.k; f++ {
		yHat += 0.5 * (math.Pow(v[f], 2) - vv[f])
	}

	return sigmoid(yHat)
}
