//-----------------------------------------------------------------------------
/*

Minimum Phase Bandwidth Limited Steps

See:
https://www.experimentalscene.com/articles/minbleps.php
https://www.cs.cmu.edu/~eli/papers/icmc01-hardsync.pdf

*/
//-----------------------------------------------------------------------------

package core

import (
	"math"
	"math/cmplx"
)

//-----------------------------------------------------------------------------

// Sinc function (pi * x variant).
func Sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	x *= Pi
	return math.Sin(x) / x
}

// BlackmanWindow returns a Blackman window with n elements.
func BlackmanWindow(n int) []float64 {
	w := make([]float64, n)
	m := float64(n - 1)
	if n == 1 {
		w[0] = 1
	} else {
		for i := 0; i < n; i++ {
			f1 := Tau * float64(i) / m
			f2 := 2 * f1
			w[i] = 0.42 - (0.5 * math.Cos(f1)) + (0.08 * math.Cos(f2))
		}
	}
	return w
}

// RealCepstrum returns the real cepstrum of a real signal.
func RealCepstrum(signal []float64) []float64 {
	freq := DFT(toComplex128(signal))
	// calculate the log of the absolute value
	for i := range freq {
		freq[i] = complex(math.Log(cmplx.Abs(freq[i])), 0)
	}
	// back to time domain
	time := InverseDFT(freq)
	// output the real part
	return toFloat64(time)
}

// MinimumPhase returns the minimum phase reconstruction of a signal.
func MinimumPhase(realCepstrum []float64) []float64 {
	n := len(realCepstrum)
	nd2 := n / 2
	realTime := make([]float64, n)
	if (n % 2) == 1 {
		realTime[0] = realCepstrum[0]
		for i := 1; i < nd2; i++ {
			realTime[i] = 2 * realCepstrum[i]
		}
		for i := nd2; i < n; i++ {
			realTime[i] = 0
		}
	} else {
		realTime[0] = realCepstrum[0]
		for i := 1; i < nd2; i++ {
			realTime[i] = 2 * realCepstrum[i]
		}
		realTime[nd2] = realCepstrum[nd2]
		for i := nd2 + 1; i < n; i++ {
			realTime[i] = 0
		}
	}
	freq := DFT(toComplex128(realTime))
	for i := range freq {
		freq[i] = cmplx.Exp(freq[i])
	}
	time := InverseDFT(freq)
	return toFloat64(time)
}

// GenerateMinBLEP returns a minimum phase bandwidth limited step.
func GenerateMinBLEP(zeroCrossings, overSampling int) []float64 {
	n := (2 * zeroCrossings * overSampling) + 1
	// generate sinc
	sinc := make([]float64, n)
	k := 1.0 / float64(overSampling)
	for i := 0; i < n; i++ {
		sinc[i] = Sinc(k*float64(i) - float64(zeroCrossings))
	}
	// window the sinc
	window := BlackmanWindow(n)
	for i := 0; i < n; i++ {
		sinc[i] *= window[i]
	}
	// minimum phase reconstruction
	realCepstrum := RealCepstrum(sinc)
	minPhase := MinimumPhase(realCepstrum)
	// integrate into minBLEP
	minBLEP := make([]float64, n)
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += minPhase[i]
		minBLEP[i] = sum
	}
	// Normalize
	scale := 1.0 / minBLEP[n-1]
	for i := 0; i < n; i++ {
		minBLEP[i] *= scale
	}
	return minBLEP
}

//-----------------------------------------------------------------------------
