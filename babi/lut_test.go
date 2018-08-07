//-----------------------------------------------------------------------------
/*

Lookup Table Testing

*/
//-----------------------------------------------------------------------------

package babi

import (
	"math"
	"testing"
)

//-----------------------------------------------------------------------------

const TEST_SIZE = 10007 // prime to give the LERP a workout
const TEST_LIMIT = 5 * TEST_SIZE
const MAX_ERR = 1e-5 // 1 part in 100000 - should be fine for 16 bit sound

func Test_Cos(t *testing.T) {
	dx := TAU / TEST_SIZE
	for i := -TEST_LIMIT; i < TEST_LIMIT; i++ {
		x := float64(i) * dx
		y0 := float64(Cos(float32(x)))
		y1 := math.Cos(x)
		err := math.Abs(y0 - y1)
		if err >= MAX_ERR {
			t.Logf("i %d x %e y0 %e y1 %e err %e", i, x, y0, y1, err)
			t.Error("FAIL")
		}
	}
}

func Test_Sin(t *testing.T) {
	dx := TAU / TEST_SIZE
	for i := -TEST_LIMIT; i < TEST_LIMIT; i++ {
		x := float64(i) * dx
		y0 := float64(Sin(float32(x)))
		y1 := math.Sin(x)
		err := math.Abs(y0 - y1)
		if err >= MAX_ERR {
			t.Logf("i %d x %e y0 %e y1 %e err %e", i, x, y0, y1, err)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

// LUT based babi.Cos
func benchmark_babi_Cos(theta float32, b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		Cos(theta)
	}
}

func Benchmark_babi_Cos0(b *testing.B)   { benchmark_babi_Cos(0, b) }
func Benchmark_babi_Cos1(b *testing.B)   { benchmark_babi_Cos(1, b) }
func Benchmark_babi_Cos10(b *testing.B)  { benchmark_babi_Cos(10, b) }
func Benchmark_babi_Cos100(b *testing.B) { benchmark_babi_Cos(100, b) }

// standard math.Cos
func benchmark_math_Cos(theta float32, b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		math.Cos(float64(theta))
	}
}

func Benchmark_math_Cos0(b *testing.B)   { benchmark_math_Cos(0, b) }
func Benchmark_math_Cos1(b *testing.B)   { benchmark_math_Cos(1, b) }
func Benchmark_math_Cos10(b *testing.B)  { benchmark_math_Cos(10, b) }
func Benchmark_math_Cos100(b *testing.B) { benchmark_math_Cos(100, b) }

//-----------------------------------------------------------------------------
