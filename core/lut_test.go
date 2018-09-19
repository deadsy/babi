//-----------------------------------------------------------------------------
/*

Lookup Table Testing

*/
//-----------------------------------------------------------------------------

package core

import (
	"math"
	"testing"
)

//-----------------------------------------------------------------------------

const testSize = 10007 // prime to give the LERP a workout
const testLimit = 5 * testSize
const maxCosError = 1e-5 // 1 part in 100000 - should be fine for 16 bit sound

func TestCos(t *testing.T) {
	dx := Tau / testSize
	for i := -testLimit; i < testLimit; i++ {
		x := float64(i) * dx
		y0 := float64(Cos(float32(x)))
		y1 := math.Cos(x)
		err := math.Abs(y0 - y1)
		if err >= maxCosError {
			t.Logf("i %d x %e y0 %e y1 %e err %e", i, x, y0, y1, err)
			t.Error("FAIL")
		}
	}
}

func TestSin(t *testing.T) {
	dx := Tau / testSize
	for i := -testLimit; i < testLimit; i++ {
		x := float64(i) * dx
		y0 := float64(Sin(float32(x)))
		y1 := math.Sin(x)
		err := math.Abs(y0 - y1)
		if err >= maxCosError {
			t.Logf("i %d x %e y0 %e y1 %e err %e", i, x, y0, y1, err)
			t.Error("FAIL")
		}
	}
}

// benchmarking

// LUT based babi.Cos
func benchmarkBabiCos(theta float32, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Cos(theta)
	}
}

func BenchmarkBabiCos0(b *testing.B)   { benchmarkBabiCos(0, b) }
func BenchmarkBabiCos1(b *testing.B)   { benchmarkBabiCos(1, b) }
func BenchmarkBabiCos10(b *testing.B)  { benchmarkBabiCos(10, b) }
func BenchmarkBabiCos100(b *testing.B) { benchmarkBabiCos(100, b) }

// standard math.Cos
func benchmarkMathCos(theta float32, b *testing.B) {
	for n := 0; n < b.N; n++ {
		math.Cos(float64(theta))
	}
}

func BenchmarkMathCos0(b *testing.B)   { benchmarkMathCos(0, b) }
func BenchmarkMathCos1(b *testing.B)   { benchmarkMathCos(1, b) }
func BenchmarkMathCos10(b *testing.B)  { benchmarkMathCos(10, b) }
func BenchmarkMathCos100(b *testing.B) { benchmarkMathCos(100, b) }

//-----------------------------------------------------------------------------

const maxPowError = 5e-5

func TestPow2(t *testing.T) {
	dx := 1.0 / testSize
	for i := -testLimit; i < testLimit; i++ {
		x := float64(i) * dx
		y0 := float64(Pow2(float32(x)))
		y1 := math.Pow(2, x)
		err := math.Abs(y0-y1) / y1
		if err >= maxPowError {
			t.Logf("i %d x %e y0 %e y1 %e err %e", i, x, y0, y1, err)
			t.Error("FAIL")
		}
	}
}

// benchmarking

// LUT based babi.Pow2
func benchmarkBabiPow2(x float32, b *testing.B) {
	for n := 0; n < b.N; n++ {
		Pow2(x)
	}
}

func BenchmarkBabiPow2_0(b *testing.B)  { benchmarkBabiPow2(0.0, b) }
func BenchmarkBabiPow2_11(b *testing.B) { benchmarkBabiPow2(1.1, b) }
func BenchmarkBabiPow2_22(b *testing.B) { benchmarkBabiPow2(2.2, b) }
func BenchmarkBabiPow2_73(b *testing.B) { benchmarkBabiPow2(7.3, b) }

// standard math.Pow
func benchmarkMathPow2(x float64, b *testing.B) {
	for n := 0; n < b.N; n++ {
		math.Pow(2, x)
	}
}

func BenchmarkMathPow2_0(b *testing.B)  { benchmarkMathPow2(0.0, b) }
func BenchmarkMathPow2_11(b *testing.B) { benchmarkMathPow2(1.1, b) }
func BenchmarkMathPow2_22(b *testing.B) { benchmarkMathPow2(2.2, b) }
func BenchmarkMathPow2_73(b *testing.B) { benchmarkMathPow2(7.3, b) }

//-----------------------------------------------------------------------------
