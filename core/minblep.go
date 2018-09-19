//-----------------------------------------------------------------------------
/*

Minimum Phase Bandwidth Limited Steps

See:
https://www.experimentalscene.com/articles/minbleps.php
https://www.cs.cmu.edu/~eli/papers/icmc01-hardsync.pdf

*/
//-----------------------------------------------------------------------------

package core

import "math"

//-----------------------------------------------------------------------------

// Sinc Function
func Sinc(x float64) float64 {
	if x == 0 {
		return 1
	}
	px := Pi * x
	return math.Sin(px) / px
}

// BlackmanWindow returns a Blackman window with n elements
func BlackmanWindow(n int) []float64 {
	w := make([]float64, n)
	m := float64(n - 1)
	if n == 1 {
		w[0] = 1
	} else {
		for i := 0; i < n; i++ {
			f1 := 2 * Pi * float64(i) / m
			f2 := 2 * f1
			w[i] = 0.42 - (0.5 * math.Cos(f1)) + (0.08 * math.Cos(f2))
		}
	}
	return w
}

// DFT returns the discrete fourier transform of the complex input.
func DFT(realTime, imagTime []float64) (realFreq, imagFreq []float64) {
	n := len(realTime)
	nInv := 1.0 / float64(n)
	realFreq = make([]float64, n)
	imagFreq = make([]float64, n)
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			p := 2 * Pi * float64(k*i) * nInv
			sr := math.Cos(p)
			si := -math.Sin(p)
			realFreq[k] += (realTime[i] * sr) - (imagTime[i] * si)
			imagFreq[k] += (realTime[i] * si) + (imagTime[i] * sr)
		}
	}
	return
}

// InverseDFT returns the inverse discrete fourier transform of the complex input.
func InverseDFT(realFreq, imagFreq []float64) (realTime, imagTime []float64) {
	n := len(realFreq)
	nInv := 1.0 / float64(n)
	realTime = make([]float64, n)
	imagTime = make([]float64, n)
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			p := 2 * Pi * float64(k*i) * nInv
			sr := math.Cos(p)
			si := -math.Sin(p)
			realTime[k] += (realFreq[i] * sr) + (imagFreq[i] * si)
			imagTime[k] += (realFreq[i] * si) - (imagFreq[i] * sr)
		}
		realTime[k] *= nInv
		imagTime[k] *= nInv
	}
	return
}

// Cabs returns the complex absolute value.
func Cabs(r, i float64) float64 {
	return math.Sqrt((r * r) + (i * i))
}

// Cexp returns the complex exponential.
func Cexp(r, i float64) (zr, zi float64) {
	er := math.Exp(r)
	zr = er * math.Cos(i)
	zi = er * math.Sin(i)
	return
}

// RealCepstrum computes the real cepstrum of a real signal.
func RealCepstrum(signal []float64) []float64 {
	n := len(signal)
	// create complex time domain input
	realTime := make([]float64, n)
	imagTime := make([]float64, n)
	copy(realTime, signal)
	// convert to frequency domain
	realFreq, imagFreq := DFT(realTime, imagTime)
	// calculate the log of the absolute value
	for i := 0; i < n; i++ {
		realFreq[i] = math.Log(Cabs(realFreq[i], imagFreq[i]))
		imagFreq[i] = 0
	}
	// back to time domain
	realTime, imagTime = InverseDFT(realFreq, imagFreq)
	// output the real part
	return realTime
}

// MinimumPhase computes the minimum phase reconstruction of a signal.
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
	imagTime := make([]float64, n)
	realFreq, imagFreq := DFT(realTime, imagTime)
	for i := 0; i < n; i++ {
		realFreq[i], imagFreq[i] = Cexp(realFreq[i], imagFreq[i])
	}
	realTime, _ = InverseDFT(realFreq, imagFreq)
	return realTime
}

// GenerateMinBLEP returns a MinBLEP as an array of floats.
func GenerateMinBLEP(zeroCrossings, overSampling int) []float64 {
	n := (zeroCrossings * 2 * overSampling) + 1
	// generate sinc
	buffer1 := make([]float64, n)
	a := float64(-zeroCrossings)
	b := float64(zeroCrossings)
	for i := 0; i < n; i++ {
		r := float64(i) / float64(n-1)
		buffer1[i] = Sinc(a + (r * (b - a)))
	}
	// window the sinc
	buffer2 := BlackmanWindow(n)
	for i := 0; i < n; i++ {
		buffer1[i] *= buffer2[i]
	}
	// minimum phase reconstruction
	buffer2 = RealCepstrum(buffer1)
	buffer1 = MinimumPhase(buffer2)
	// integrate into minBLEP
	minBLEP := make([]float64, n)
	a = 0
	for i := 0; i < n; i++ {
		a += buffer1[i]
		minBLEP[i] = a
	}
	// Normalize
	a = minBLEP[n-1]
	a = 1.0 / a
	for i := 0; i < n; i++ {
		minBLEP[i] *= a
	}
	return minBLEP
}

//-----------------------------------------------------------------------------
