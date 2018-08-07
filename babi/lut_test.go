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
const MAX_COS_ERR = 1e-5 // 1 part in 100000 - should be fine for 16 bit sound

func Test_Cos(t *testing.T) {
	dx := TAU / TEST_SIZE
	for i := -TEST_LIMIT; i < TEST_LIMIT; i++ {
		x := float64(i) * dx
		y0 := float64(Cos(float32(x)))
		y1 := math.Cos(x)
		err := math.Abs(y0 - y1)
		if err >= MAX_COS_ERR {
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
		if err >= MAX_COS_ERR {
			t.Logf("i %d x %e y0 %e y1 %e err %e", i, x, y0, y1, err)
			t.Error("FAIL")
		}
	}
}

// benchmarking

// LUT based babi.Cos
func benchmark_babi_Cos(theta float32, b *testing.B) {
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
	for n := 0; n < b.N; n++ {
		math.Cos(float64(theta))
	}
}

func Benchmark_math_Cos0(b *testing.B)   { benchmark_math_Cos(0, b) }
func Benchmark_math_Cos1(b *testing.B)   { benchmark_math_Cos(1, b) }
func Benchmark_math_Cos10(b *testing.B)  { benchmark_math_Cos(10, b) }
func Benchmark_math_Cos100(b *testing.B) { benchmark_math_Cos(100, b) }

//-----------------------------------------------------------------------------

const MAX_POW_ERR = 5e-5

func Test_Pow2(t *testing.T) {
	dx := 1.0 / TEST_SIZE
	for i := -TEST_LIMIT; i < TEST_LIMIT; i++ {
		x := float64(i) * dx
		y0 := float64(Pow2(float32(x)))
		y1 := math.Pow(2, x)
		err := math.Abs(y0-y1) / y1
		if err >= MAX_POW_ERR {
			t.Logf("i %d x %e y0 %e y1 %e err %e", i, x, y0, y1, err)
			t.Error("FAIL")
		}
	}
}

// benchmarking

// LUT based babi.Pow2
func benchmark_babi_Pow2(x float32, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Pow2(x)
	}
}

func Benchmark_babi_Pow2_0(b *testing.B)  { benchmark_babi_Pow2(0.0, b) }
func Benchmark_babi_Pow2_11(b *testing.B) { benchmark_babi_Pow2(1.1, b) }
func Benchmark_babi_Pow2_22(b *testing.B) { benchmark_babi_Pow2(2.2, b) }
func Benchmark_babi_Pow2_73(b *testing.B) { benchmark_babi_Pow2(7.3, b) }

// standard math.Pow
func benchmark_math_Pow2(x float64, b *testing.B) {
	for n := 0; n < b.N; n++ {
		math.Pow(2, x)
	}
}

func Benchmark_math_Pow2_0(b *testing.B)  { benchmark_math_Pow2(0.0, b) }
func Benchmark_math_Pow2_11(b *testing.B) { benchmark_math_Pow2(1.1, b) }
func Benchmark_math_Pow2_22(b *testing.B) { benchmark_math_Pow2(2.2, b) }
func Benchmark_math_Pow2_73(b *testing.B) { benchmark_math_Pow2(7.3, b) }

//-----------------------------------------------------------------------------
