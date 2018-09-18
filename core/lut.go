//-----------------------------------------------------------------------------
/*

Lookup Table Based Math Functions

Faster than standard math package functions, but less accurate.

*/
//-----------------------------------------------------------------------------

package core

import "math"

//-----------------------------------------------------------------------------

func init() {
	cosLUTInit()
	powLUTInit()
}

//-----------------------------------------------------------------------------
// Cosine Lookup

const cosLUTBits = 10
const cosLUTSize = 1 << cosLUTBits
const cosFracBits = 32 - cosLUTBits
const cosFracMask = (1 << cosFracBits) - 1

var cosLUTy [cosLUTSize]float32
var cosLUTdy [cosLUTSize]float32

// cosLUTInit creates y/dy cosine lookup tables for TAU radians.
func cosLUTInit() {
	dx := Tau / cosLUTSize
	for i := 0; i < cosLUTSize; i++ {
		y0 := math.Cos(float64(i) * dx)
		y1 := math.Cos(float64(i+1) * dx)
		cosLUTy[i] = float32(y0)
		cosLUTdy[i] = float32((y1 - y0) / (1 << cosFracBits))
	}
}

// CosLookup returns the cosine of x (32 bit unsigned phase value).
func CosLookup(x uint32) float32 {
	idx := x >> cosFracBits
	return cosLUTy[idx] + float32(x&cosFracMask)*cosLUTdy[idx]
}

// Cos returns the cosine of x (radians).
func Cos(x float32) float32 {
	xi := uint32(Abs(x) * PhaseScale)
	return CosLookup(xi)
}

// Sin returns the sine of x (radians).
func Sin(x float32) float32 {
	return Cos((Pi / 2) - x)
}

// Tan returns the tangent of x (radians).
func Tan(x float32) float32 {
	return Sin(x) / Cos(x)
}

//-----------------------------------------------------------------------------
// Power Function

const powLUTBits = 7
const powLUTSize = 1 << powLUTBits
const powLUTMask = powLUTSize - 1

var powLUT0 [powLUTSize]float32
var powLUT1 [powLUTSize]float32

// powLUTInit creates the power lookup tables.
func powLUTInit() {
	for i := 0; i < powLUTSize; i++ {
		x := float64(i) / powLUTSize
		powLUT0[i] = float32(math.Pow(2, x))
		x = float64(i) / (powLUTSize * powLUTSize)
		powLUT1[i] = float32(math.Pow(2, x))
	}
}

// pow2_int returns 2 to the x where x is an integer [-126,127]
func pow2Int(x int) float32 {
	return math.Float32frombits((127 + uint32(x)) << 23)
}

// pow2_frac returns 2 to the x where x is a fraction [0,1)
func pow2Frac(x float32) float32 {
	n := int(x * (1 << (powLUTBits * 2)))
	x0 := powLUT0[(n>>powLUTBits)&powLUTMask]
	x1 := powLUT1[n&powLUTMask]
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
		nf--
		ff++
	}
	return pow2Int(nf) * pow2Frac(ff)
}

const logE2 = 1.4426950408889634 // 1.0 / math.log(2)

// PowE returns e to the x.
func PowE(x float32) float32 {
	return Pow2(logE2 * x)
}

//-----------------------------------------------------------------------------
