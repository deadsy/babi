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
	"math/bits"
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

// isPowerOf2 return true if n is a power of 2.
func isPowerOf2(x int) bool {
	return x != 0 && (x&-x) == x
}

// bitReverse reverses the first n bits of x.
func bitReverse(x, n int) int {
	return int(bits.Reverse(uint(x)) >> (bits.UintSize - uint(n)))
}

// log2 returns log base 2 of x (assumes x is a power of 2).
func log2(x int) int {
	return bits.TrailingZeros(uint(x))
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

/*

// FFT returns the (fast) discrete fourier transform of the complex input.
func FFT(in []complex128) []complex128 {
	// check input length
	n := len(in)
	if !isPowerOf2(n) {
		panic("input length is not a power of 2")
	}
	// reverse the input order
	out := make([]complex128, n)
	nbits := log2(n)
	for i := range out {
		out[i] = in[bitReverse(i, nbits)]
	}
	// run the butterflies
	kmax := 1
	for {
		if kmax >= n {
			return out
		}
		istep := kmax * 2
		for k := 0; k < kmax; k++ {
			theta := -Pi * float64(k) / float64(kmax)
			s, c := math.Sincos(theta)
			cs := complex(c, s)
			for i := k; i < n; i += istep {
				j := i + kmax
				temp := out[j] * cs
				out[j] = out[i] - temp
				out[i] = out[i] + temp
			}
		}
		kmax = istep
	}
}

*/

// InverseFFT returns the (fast) inverse discrete fourier transform of the complex input.
func InverseFFT(in []complex128) []complex128 {
	n := len(in)
	nInv := complex(1.0/float64(n), 0)
	out := make([]complex128, n)
	for i := range out {
		out[i] = cmplx.Conj(in[i])
	}
	out = FFT(out)
	for i := range out {
		out[i] = cmplx.Conj(out[i])
		out[i] *= nInv
	}
	return out
}

//-----------------------------------------------------------------------------

// fftConst contains pre-calculated fft constants.
type fftConst struct {
	n       int          // length of fft input
	reverse []int        // input reversing indices
	twiddle []complex128 // twiddle factors
}

// fftCache is a cache of pre-calculated fft constants.
var fftCache map[int]*fftConst

// fftLookup returns the fft constants for a given input length.
func fftLookup(n int) *fftConst {

	// has the cache been created?
	if fftCache == nil {
		fftCache = make(map[int]*fftConst)
	}

	// do we have the entry in the cache?
	if k, ok := fftCache[n]; ok {
		return k
	}

	// check length
	if !isPowerOf2(n) {
		panic("input length is not a power of 2")
	}
	if n < 4 {
		panic("input length has to be >= 4")
	}

	// create the entry
	k := &fftConst{}
	k.n = n

	// create the reverse indices
	k.reverse = make([]int, n)
	nbits := log2(n)
	for i := range k.reverse {
		k.reverse[i] = bitReverse(i, nbits)
	}

	// create the twiddle factors (use quadrant symmetry)
	k.twiddle = make([]complex128, n)
	nInv := 1.0 / float64(n)
	for i := 0; i < n>>2; i++ {
		theta := -Tau * float64(i) * nInv
		s, c := math.Sincos(theta)
		k.twiddle[i] = complex(c, s)
		k.twiddle[i+(n>>2)] = complex(-s, c)
		k.twiddle[i+(n>>1)] = complex(-c, -s)
		k.twiddle[i+(n>>2)+(n>>1)] = complex(s, -c)
	}

	// add it to the cache and return
	fftCache[n] = k
	return k
}

// FFTx returns the (fast) discrete fourier transform of the complex input.
func FFT(in []complex128) []complex128 {

	n := len(in)

	fk := fftLookup(n)
	// reverse the input order
	out := make([]complex128, n)
	for i := range out {
		out[i] = in[fk.reverse[i]]
	}

	// run the butterflies
	kmax := 1
	for {
		if kmax >= n {
			return out
		}
		istep := kmax * 2
		for k := 0; k < kmax; k++ {
			theta := -Pi * float64(k) / float64(kmax)
			s, c := math.Sincos(theta)
			cs := complex(c, s)
			for i := k; i < n; i += istep {
				j := i + kmax
				temp := out[j] * cs
				out[j] = out[i] - temp
				out[i] = out[i] + temp
			}
		}
		kmax = istep
	}
}

//-----------------------------------------------------------------------------
