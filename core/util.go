//-----------------------------------------------------------------------------
/*

Common Utility Functions

*/
//-----------------------------------------------------------------------------

package core

import (
	"fmt"
	"strings"
)

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

// Clamp clamps x between a and b (float32).
func Clamp(x, a, b float32) float32 {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

// ClampLo clamps x to >= a (float32).
func ClampLo(x, a float32) float32 {
	if x < a {
		return a
	}
	return x
}

// ClampInt clamps x between a and b (integer).
func ClampInt(x, a, b int) int {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

//-----------------------------------------------------------------------------

// MapLin returns a linear mapping from x = 0..1 to y = a..b.
func MapLin(x, a, b float32) float32 {
	return ((b - a) * x) + a
}

/*

// MapLog returns a logarithmic mapping from x = 0..1 to y = a..b.
func MapLog(x, a, b float32) float32 {
	if x == 0 {
		return a
	}
	if x == 1 {
		return b
	}
	panic("todo")
}

// MapExp returns an exponential mapping from x = 0..1 to y = a..b.
func MapExp(x, a, b float32) float32 {
	if x == 0 {
		return a
	}
	if x == 1 {
		return b
	}
	k := float32(math.Log2(float64(b / a)))
	return a * Pow2(k*x)
}

*/

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

// TableString return a string for a table of row by column strings.
// Each column string will be left justified and aligned.
func TableString(
	rows [][]string, // table rows [[col0, col1, col2...,colN]...]
	csize []int, // minimum column widths
	cmargin int, // column to column margin
) string {
	// how many rows?
	nrows := len(rows)
	if nrows == 0 {
		return ""
	}
	// how many columns?
	ncols := len(rows[0])
	// make sure we have a well formed csize
	if csize == nil {
		csize = make([]int, ncols)
	} else {
		if len(csize) != ncols {
			panic("len(csize) != ncols")
		}
	}
	// check that the number of columns for each row is consistent
	for i := range rows {
		if len(rows[i]) != ncols {
			panic(fmt.Sprintf("ncols row%d != ncols row0", i))
		}
	}
	// go through the strings and bump up csize widths if required
	for i := 0; i < nrows; i++ {
		for j := 0; j < ncols; j++ {
			width := len(rows[i][j])
			if (width + cmargin) >= csize[j] {
				csize[j] = width + cmargin
			}
		}
	}
	// build the row format string
	fmtCol := make([]string, ncols)
	for i, n := range csize {
		fmtCol[i] = fmt.Sprintf("%%-%ds", n)
	}
	fmtRow := strings.Join(fmtCol, "")
	// generate the row strings
	row := make([]string, nrows)
	for i, l := range rows {
		// convert []string to []interface{}
		x := make([]interface{}, len(l))
		for j, v := range l {
			x[j] = v
		}
		row[i] = fmt.Sprintf(fmtRow, x...)
	}
	// return rows and columns
	return strings.Join(row, "\n")
}

//-----------------------------------------------------------------------------

// BoolToString returns one of two strings based on the boolean.
func BoolToString(val bool, str []string) string {
	var x int
	if val {
		x = 1
	}
	return str[x]
}

//-----------------------------------------------------------------------------
