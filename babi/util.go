//-----------------------------------------------------------------------------
/*

Common Utility Functions

*/
//-----------------------------------------------------------------------------

package babi

//-----------------------------------------------------------------------------

// Clamp x between a and b
func Clamp(x, a, b float32) float32 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

// Clamp x to >= a
func ClampLo(x, a float32) float32 {
	if x < a {
		return a
	}
	return x
}

// Clamp x to <= a
func ClampHi(x, a float32) float32 {
	if x > a {
		return a
	}
	return x
}

// Linear mapping of x = 0..1 to y = a..b
func Map(x, a, b float32) float32 {
	return ((b - a) * x) + a
}

//-----------------------------------------------------------------------------
