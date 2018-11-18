//-----------------------------------------------------------------------------
/*

Common Utility Functions

*/
//-----------------------------------------------------------------------------

package core

//-----------------------------------------------------------------------------

// Min returns the minimum of two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

// Abs return the absolute value of x.
func Abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}

//-----------------------------------------------------------------------------

// Clamp clamps x between a and b.
func Clamp(x, a, b float32) float32 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

// ClampLo clamps x to >= a.
func ClampLo(x, a float32) float32 {
	if x < a {
		return a
	}
	return x
}

// ClampHi clamps x to <= a.
func ClampHi(x, a float32) float32 {
	if x > a {
		return a
	}
	return x
}

//-----------------------------------------------------------------------------

// Map returns a linear mapping from x = 0..1 to y = a..b.
func Map(x, a, b float32) float32 {
	return ((b - a) * x) + a
}

//-----------------------------------------------------------------------------

// InRange returns true if a <= x <= b.
func InRange(x, a, b float32) bool {
	return x >= a && x <= b
}

// InEnum returns true if x is in [0:max)
func InEnum(x, max int) bool {
	return x >= 0 && x < max
}

//-----------------------------------------------------------------------------

// DtoR converts degrees to radians.
func DtoR(degrees float64) float64 {
	return (Pi / 180.0) * degrees
}

// RtoD converts radians to degrees.
func RtoD(radians float64) float64 {
	return (180.0 / Pi) * radians
}

//-----------------------------------------------------------------------------

// SignExtend sign extends an n bit value to a signed integer.
func SignExtend(x int, n uint) int {
	y := uint(x) & ((1 << n) - 1)
	mask := uint(1 << (n - 1))
	return int((y ^ mask) - mask)
}

//-----------------------------------------------------------------------------
