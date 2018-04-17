package line

import "math"

// FTRLProximalClassifier implements the "Follow The Regularized Leader
// Proximal" algorithm described in
// http://www.eecs.tufts.edu/~dsculley/papers/ad-click-prediction.pdf.
type FTRLProximalClassifier struct {
	alpha float64
	beta  float64
	l1    float64
	l2    float64

	n Vector // Squared sum of past gradients
	z Vector // Weights
	w Vector // Lazy weights
}

// NewFTRLProximalClassifier instantiates and returns a *FTRLProximalClassifier.
func NewFTRLProximalClassifier(alpha, beta, l1, l2 float64) *FTRLProximalClassifier {
	return &FTRLProximalClassifier{
		alpha: alpha,
		beta:  beta,
		l1:    l1,
		l2:    l2,
		n:     make(Vector),
		z:     make(Vector),
	}
}

// FitPartial of FTRLProximalClassifier.
func (ftrl *FTRLProximalClassifier) FitPartial(x Vector, y float64) (yPred float64) {
	yPred = ftrl.PredictPartial(x)
	loss := yPred - y
	for i, xi := range x {
		g := loss * xi // Gradient of loss w.r.t wi
		g2 := g * g
		sigma := (math.Sqrt(ftrl.n[i]+g2) - math.Sqrt(ftrl.n[i])) / ftrl.alpha
		ftrl.z[i] += g - sigma*ftrl.w[i]
		ftrl.n[i] += g2
	}
	return
}

// PredictPartial of FTRLProximalClassifier.
func (ftrl *FTRLProximalClassifier) PredictPartial(x Vector) (yPred float64) {
	var (
		w   = make(Vector)
		wTx float64 // Inner product
	)
	for i, xi := range x {
		z, ok := ftrl.z[i]
		if !ok {
			ftrl.z[i] = 0
			ftrl.n[i] = 0
		}
		sign := 1.0
		if z < 0 {
			sign = -1.0
		}
		if sign*z <= ftrl.l1 {
			w[i] = 0
		} else {
			w[i] = (sign*ftrl.l1 - z) / ((ftrl.beta+math.Sqrt(ftrl.n[i]))/ftrl.alpha + ftrl.l2)
			wTx += w[i] * xi
		}
	}
	ftrl.w = w
	return sigmoid(wTx)
}
