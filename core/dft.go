//-----------------------------------------------------------------------------
/*

Discrete Fourier Transform

See:
https://en.wikipedia.org/wiki/Discrete_Fourier_transform
https://github.com/takatoh/fft

*/
//-----------------------------------------------------------------------------

package core

import (
	"math"
	"math/cmplx"
)

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

// FFT returns the (fast) discrete fourier transform of the complex input.
func FFT(in []complex128) []complex128 {
	n := len(in)
	out := make([]complex128, n)
	copy(out, in)
	j := 0
	for i := 0; i < n; i++ {
		if i < j {
			out[i], out[j] = out[j], out[i]
		}
		m := n / 2
		for {
			if j < m {
				break
			}
			j = j - m
			m = m / 2
			if m < 2 {
				break
			}
		}
		j = j + m
	}
	kmax := 1
	for {
		if kmax >= n {
			return out
		}
		istep := kmax * 2
		for k := 0; k < kmax; k++ {
			theta := complex(0, -Pi*float64(k)/float64(kmax))
			for i := k; i < n; i += istep {
				j := i + kmax
				temp := out[j] * cmplx.Exp(theta)
				out[j] = out[i] - temp
				out[i] = out[i] + temp
			}
		}
		kmax = istep
	}
}

// InverseFFT returns the (fast) inverse discrete fourier transform of the complex input.
func InverseFFT(in []complex128) []complex128 {
	out := make([]complex128, len(in))
	for i := range out {
		out[i] = complex(real(in[i]), -imag(in[i]))
	}
	out = FFT(out)
	for i := range out {
		out[i] = complex(real(out[i]), -imag(out[i]))
	}
	return out
}

//-----------------------------------------------------------------------------
