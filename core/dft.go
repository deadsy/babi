//-----------------------------------------------------------------------------
/*

Discrete Fourier Transform

See:
https://en.wikipedia.org/wiki/Discrete_Fourier_transform

*/
//-----------------------------------------------------------------------------

package core

import "math"

//-----------------------------------------------------------------------------

// toComplex128 converts a slice of float values to complex values.
// The imaginary part is set to zero.
func toComplex128(in []float64) []complex128 {
	out := make([]complex128, len(in))
	for i := range out {
		out[i] = complex(in[i], 0)
	}
	return out
}

// toFloat64 converts a slice of complex values to float values by taking the real part.
func toFloat64(in []complex128) []float64 {
	out := make([]float64, len(in))
	for i := range out {
		out[i] = real(in[i])
	}
	return out
}

//-----------------------------------------------------------------------------

// DFT returns the discrete fourier transform of the complex input.
func DFT(in []complex128) []complex128 {
	n := len(in)
	nInv := 1.0 / float64(n)
	out := make([]complex128, n)
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			p := -Tau * float64(k*i) * nInv
			s, c := math.Sincos(p)
			out[k] += in[i] * complex(c, s)
		}
	}
	return out
}

// InverseDFT returns the inverse discrete fourier transform of the complex input.
func InverseDFT(in []complex128) []complex128 {
	n := len(in)
	nInv := 1.0 / float64(n)
	out := make([]complex128, n)
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			p := Tau * float64(k*i) * nInv
			s, c := math.Sincos(p)
			out[k] += in[i] * complex(c, s)
		}
		out[k] *= complex(nInv, 0)
	}
	return out
}

//-----------------------------------------------------------------------------
