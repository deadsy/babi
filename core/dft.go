//-----------------------------------------------------------------------------
/*

Discrete Fourier Transform

*/
//-----------------------------------------------------------------------------

package core

import "math"

//-----------------------------------------------------------------------------

/*

func DFT(in []complex128) out []complex128 {

	n := len(in)

	nInv := 1.0 / float64(n)

	out = make([]complex128, n)

	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {

			p := Tau * float64(k*i) * nInv

      s, c := math.Sincos(p)


			sr := math.Cos(p)
			si := -math.Sin(p)


			outRe[k] += (inRe[i] * c) + (inIm[i] * s)
			outIm[k] += (inRe[i] * -s) + (inIm[i] * c)
		}
	}

  return out
}

*/

//-----------------------------------------------------------------------------

// DFT returns the discrete fourier transform of the complex input.
func DFT(inRe, inIm []float64) (outRe, outIm []float64) {
	n := len(inRe)
	nInv := 1.0 / float64(n)
	outRe = make([]float64, n)
	outIm = make([]float64, n)
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			p := 2 * Pi * float64(k*i) * nInv
			sr := math.Cos(p)
			si := -math.Sin(p)
			outRe[k] += (inRe[i] * sr) - (inIm[i] * si)
			outIm[k] += (inRe[i] * si) + (inIm[i] * sr)
		}
	}
	return
}

// InverseDFT returns the inverse discrete fourier transform of the complex input.
func InverseDFT(inRe, inIm []float64) (outRe, outIm []float64) {
	n := len(inRe)
	nInv := 1.0 / float64(n)
	outRe = make([]float64, n)
	outIm = make([]float64, n)
	for k := 0; k < n; k++ {
		for i := 0; i < n; i++ {
			p := 2 * Pi * float64(k*i) * nInv
			sr := math.Cos(p)
			si := -math.Sin(p)
			outRe[k] += (inRe[i] * sr) + (inIm[i] * si)
			outIm[k] += (inRe[i] * si) - (inIm[i] * sr)
		}
		outRe[k] *= nInv
		outIm[k] *= nInv
	}
	return
}

//-----------------------------------------------------------------------------
