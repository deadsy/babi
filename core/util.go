//-----------------------------------------------------------------------------
/*

Common Utility Functions

*/
//-----------------------------------------------------------------------------

package core

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

// Map returns a linear mapping from x = 0..1 to y = a..b.
func Map(x, a, b float32) float32 {
	return ((b - a) * x) + a
}

// InRange returns true if a <= x <= b.
func InRange(x, a, b float32) bool {
	return x >= a && x <= b
}

// InEnum returns true if x is in [0:max)
func InEnum(x, max int) bool {
	return x >= 0 && x < max
}

//-----------------------------------------------------------------------------
