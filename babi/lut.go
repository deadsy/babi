//-----------------------------------------------------------------------------
/*

Lookup Table Based Math Functions

Faster than standard math package functions, but less accurate.

*/
//-----------------------------------------------------------------------------

package babi

import "math"

//-----------------------------------------------------------------------------

func init() {
	cos_lut_init()
	pow_lut_init()
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
// Power Function

const POW_LUT_BITS = 7
const POW_LUT_SIZE = 1 << POW_LUT_BITS
const POW_LUT_MASK = POW_LUT_SIZE - 1

var POW_LUT0 [POW_LUT_SIZE]float32
var POW_LUT1 [POW_LUT_SIZE]float32

func pow_lut_init() {
	for i := 0; i < POW_LUT_SIZE; i++ {
		x := float64(i) / POW_LUT_SIZE
		POW_LUT0[i] = float32(math.Pow(2, x))
		x = float64(i) / (POW_LUT_SIZE * POW_LUT_SIZE)
		POW_LUT1[i] = float32(math.Pow(2, x))
	}
}

// pow2_int returns 2 to the x where x is an integer [-126,127]
func pow2_int(x int) float32 {
	return math.Float32frombits((127 + uint32(x)) << 23)
}

// pow2_frac returns 2 to the x where x is a fraction [0,1)
func pow2_frac(x float32) float32 {
	n := int(x * (1 << (POW_LUT_BITS * 2)))
	x0 := POW_LUT0[(n>>POW_LUT_BITS)&POW_LUT_MASK]
	x1 := POW_LUT1[n&POW_LUT_MASK]
	return x0 * x1
}

// Pow2 returns 2 to the x.
func Pow2(x float32) float32 {
	if x == 0 {
		return 1
	}
	nf := int(math.Trunc(float64(x)))
	ff := x - float32(nf)
	if ff < 0 {
		nf -= 1
		ff += 1
	}
	return pow2_int(nf) * pow2_frac(ff)
}

const LOG_E2 = 1.4426950408889634 // 1.0 / math.log(2)

// PowE returns e to the x.
func PowE(x float32) float32 {
	return Pow2(LOG_E2 * x)
}

//-----------------------------------------------------------------------------
