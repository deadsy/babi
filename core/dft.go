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

// fftConst contains pre-calculated fft constants.
type fftConst struct {
	hn     int          // half length
	hmask  int          // half mask for mod n/2
	stages int          // number of butterfly stages
	rev    []int        // input reversing indices
	w      []complex128 // twiddle factors
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

	// create the reverse indices
	k.rev = make([]int, n)
	nbits := log2(n)
	for i := range k.rev {
		k.rev[i] = bitReverse(i, nbits)
	}

	// create the half variables
	k.hn = n >> 1
	k.hmask = k.hn - 1

	// number of butterfly stages
	k.stages = nbits

	// create the twiddle factors
	k.w = make([]complex128, k.hn)
	nInv := 1.0 / float64(n)
	for i := range k.w {
		theta := -Tau * float64(i) * nInv
		s, c := math.Sincos(theta)
		k.w[i] = complex(c, s)
	}

	// add it to the cache and return
	fftCache[n] = k
	return k
}

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

// FFT returns the (fast) discrete fourier transform of the complex input.
func FFT(in []complex128) []complex128 {

	n := len(in)
	fk := fftLookup(n)

	// reverse the input order
	out := make([]complex128, n)
	for i := range out {
		out[i] = in[fk.rev[i]]
	}

	// run the butterflies
	kmax := 1
	mul := fk.hn
	for {
		if kmax >= n {
			break
		}
		istep := kmax * 2
		for k := 0; k < kmax; k++ {
			w := fk.w[k*mul]
			for i := k; i < n; i += istep {
				j := i + kmax
				tmp := out[j] * w
				out[j] = out[i] - tmp
				out[i] += tmp
			}
		}
		mul >>= 1
		kmax = istep
	}

	return out
}

//-----------------------------------------------------------------------------
// test code

// FFTx returns the (fast) discrete fourier transform of the complex input.
func FFTx(in []complex128) []complex128 {

	n := len(in)
	fk := fftLookup(n)

	// reverse the input order
	out := make([]complex128, n)
	for i := range out {
		out[i] = in[fk.rev[i]]
	}

	// run the butterflies
	oneMask := 1
	hiMask := -1
	loMask := 0
	shift := uint(fk.stages - 1)

	for s := 0; s < fk.stages; s++ {
		for i := 0; i < fk.hn; i++ {
			j := (i&hiMask)<<1 | (i & loMask)
			k := j | oneMask
			tmp := out[k] * fk.w[(i<<shift)&fk.hmask]
			out[k] = out[j] - tmp
			out[j] += tmp
		}
		shift--
		oneMask <<= 1
		hiMask <<= 1
		loMask = (loMask << 1) | 1
	}

	return out
}

//-----------------------------------------------------------------------------
