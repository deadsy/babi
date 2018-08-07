//-----------------------------------------------------------------------------
/*

Lookup Table Based Math Functions

*/
//-----------------------------------------------------------------------------

package babi

import "math"

//-----------------------------------------------------------------------------

func init() {
	cos_lut_init()
}

//-----------------------------------------------------------------------------
// Cosine Lookup

const COS_LUT_BITS = 10
const COS_LUT_SIZE = 1 << COS_LUT_BITS
const COS_FRAC_BITS = 32 - COS_LUT_BITS
const COS_FRAC_MASK = (1 << COS_FRAC_BITS) - 1

var COS_LUT_y [COS_LUT_SIZE]float32
var COS_LUT_dy [COS_LUT_SIZE]float32

// cos_lut_init creates y/dy cosine lookup tables for TAU radians.
func cos_lut_init() {
	dx := TAU / COS_LUT_SIZE
	for i := 0; i < COS_LUT_SIZE; i++ {
		y0 := math.Cos(float64(i) * dx)
		y1 := math.Cos(float64(i+1) * dx)
		COS_LUT_y[i] = float32(y0)
		COS_LUT_dy[i] = float32((y1 - y0) / (1 << COS_FRAC_BITS))
	}
}

// CosLookup returns the cosine of x (32 bit unsigned phase value).
func CosLookup(x uint32) float32 {
	idx := x >> COS_FRAC_BITS
	return COS_LUT_y[idx] + float32(x&COS_FRAC_MASK)*COS_LUT_dy[idx]
}

const PHASE_SCALE = (1 << 32) / TAU

// Cos returns the cosine of x (radians).
func Cos(x float32) float32 {
	xi := uint32(Abs(x) * PHASE_SCALE)
	return CosLookup(xi)
}

// Sin returns the sine of x (radians).
func Sin(x float32) float32 {
	return Cos((PI / 2) - x)
}

//-----------------------------------------------------------------------------

//-----------------------------------------------------------------------------
