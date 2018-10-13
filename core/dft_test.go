//-----------------------------------------------------------------------------
/*

Discrete Fourier Transform Testing

*/
//-----------------------------------------------------------------------------

package core

import "testing"

//-----------------------------------------------------------------------------

func TestPowerOf2(t *testing.T) {
	test := []struct {
		x      int
		result bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{4, true},
		{8, true},
		{16, true},
		{128, true},
		{129, false},
		{127, false},
		{256, true},
		{1 << 30, true},
		{1<<30 + 3, false},
	}
	for _, v := range test {
		if isPowerOf2(v.x) != v.result {
			t.Logf("for %d expected %v, actual %v\n", v.x, v.result, !v.result)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestReverse(t *testing.T) {
	test := []struct {
		x, n, y int
	}{
		{0, 3, 0},
		{1, 3, 4},
		{2, 3, 2},
		{3, 3, 6},
		{4, 3, 1},
		{5, 3, 5},
		{6, 3, 3},
		{7, 3, 7},
		{1, 32, 1 << 31},
		{5, 32, (1 << 31) | (1 << 29)},
	}
	for _, v := range test {
		if bitReverse(v.x, v.n) != v.y {
			t.Logf("for bitReverse(%d, %d) expected %d, actual %d\n", v.x, v.n, v.y, bitReverse(v.x, v.n))
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestLog2(t *testing.T) {
	test := []struct {
		x, y int
	}{
		{1, 0},
		{2, 1},
		{4, 2},
		{8, 3},
		{16, 4},
		{32, 5},
		{64, 6},
		{1024, 10},
		{1 << 31, 31},
	}
	for _, v := range test {
		if log2(v.x) != v.y {
			t.Logf("for log2(%d) expected %d, actual %d\n", v.x, v.y, log2(v.x))
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestDFT(t *testing.T) {

	realTimei := []uint64{
		0x400205ba60000000, 0xc0105145b8000000, 0x4006022400000000, 0x3fdf166b00000000, 0xbfc06d7300000000, 0x4022e93768000000, 0xc010994188000000, 0x4015b567f0000000,
		0x3fe11de500000000, 0x401597d5f0000000, 0xbfffed4580000000, 0x401f5285e0000000, 0xc01155b6d8000000, 0xc0079b4e40000000, 0x40189e32d0000000, 0x4020c2d4d0000000,
		0xc02135b4b9000000, 0x4021f91be8000000, 0x3fe0a31500000000, 0xc0208ec9c4000000, 0xc0189f75dc000000, 0x400a1dc620000000, 0x401f37f6d0000000, 0xc0082d5880000000,
		0xc0216ee2bb000000, 0xc02332f6c8000000, 0xbfeb122580000000, 0xc02179e60c000000, 0xc014f006cc000000, 0x4022d34b10000000, 0x4020169c60000000, 0x401c12d650000000,
	}
	imagTimei := []uint64{
		0xc012aaaf70000000, 0x3fe9725600000000, 0xc003f785b0000000, 0x4014d1e4d0000000, 0x3fd00b9900000000, 0x400ad5f740000000, 0x3fe43a6600000000, 0xc0226dc4ef800000,
		0xbff3f4ba60000000, 0x402145fdb0000000, 0x40213b7e08000000, 0x4011ad1a90000000, 0xc01141abc0000000, 0x4013152ee0000000, 0x4006658860000000, 0xc0075a2aa0000000,
		0x400e0ecd60000000, 0xc01ab8d958000000, 0xbff32aa360000000, 0x401e67f0c0000000, 0x401a5609e0000000, 0xc00b256410000000, 0xc015aebb58000000, 0x401f784460000000,
		0xc007f13f30000000, 0x400dddffa0000000, 0x4022423c10000000, 0x3ffc5d66c0000000, 0x40092b2c80000000, 0x401cb1b0e0000000, 0xbff35740a0000000, 0x4020f57368000000,
	}
	realFreqi := []uint64{
		0x4030e36bc8000000, 0x3ff1c43d82f0ed34, 0x4005ebbb8d1af263, 0x4024e51710b69832, 0x403eddd4d2e90096, 0xc0486a090e270dc7, 0xc0431b413c53f7d7, 0x4033f3e7cd1e83fa,
		0xc04abece6f1fffaf, 0x403925007fb3b762, 0xc03811a6de63c6ce, 0xbfc709e4c5a1f3c0, 0xc02df92c7c87e75d, 0x4042a3f6405ff870, 0xc012d0d9dd9b9e15, 0x4028235452588dec,
		0xc0448c0315fffefc, 0x40377e841518c04e, 0x3ff23680336b7c8c, 0x4013a079220dd538, 0x401913c6d45bfffb, 0x404ff87dc8225f4f, 0xc04283da1e0272c5, 0x4042ff4b24f544cf,
		0xc045ec7e76e00092, 0xc037d9871d906980, 0x404dc28e76c4d8b2, 0xc0401eb998bb2687, 0xc030e02649bc0369, 0x40368857ffde5892, 0x404ba0886809dc50, 0x4033f2602ded7844,
	}
	imagFreqi := []uint64{
		0x404a7c9216200000, 0xc040a958e3b80605, 0xc0218a1f657e26e4, 0x403ead0fea47a803, 0x403820608dacea44, 0x3ffec2a055fc7518, 0x4001c6ae5cd31480, 0x4048fa70042be868,
		0xc0351f126e000286, 0xc04dcca2d7858db1, 0x404037de84804d9c, 0x4041a4dd9a5997c8, 0x4031da6fc6491984, 0xc02e45cac13d64e8, 0xbffbbe5e60c3c740, 0xc04efe3b83bb960a,
		0xc03e3b11b44001e4, 0x4021e141fe57f797, 0xc018ac7a12a47f3e, 0xc050d10720333039, 0xc025d243b359d00e, 0xc037892873ecc5f2, 0x4043493df5e5ac1f, 0xc0205bc24271ef46,
		0x3fdca5bd7ffefdd0, 0xc00906141f82299d, 0xc0314fe5919864da, 0xc0322bafd06071c2, 0xc052955cd6924706, 0xc0318b33a194d34e, 0xc029341462b2f1ac, 0x404726a89f4558d2,
	}

	realTime := convertToFloat64(realTimei)
	imagTime := convertToFloat64(imagTimei)
	realFreq := convertToFloat64(realFreqi)
	imagFreq := convertToFloat64(imagFreqi)

	n := len(realTimei)
	time := make([]complex128, n)
	freq := make([]complex128, n)
	for i := 0; i < n; i++ {
		time[i] = complex(realTime[i], imagTime[i])
		freq[i] = complex(realFreq[i], imagFreq[i])
	}

	x := DFT(time)

	for i := 0; i < n; i++ {
		if !equal(real(x[i]), real(freq[i]), 1e-11) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, real(freq[i]), real(x[i]))
			t.Error("FAIL")
		}
		if !equal(imag(x[i]), imag(freq[i]), 1e-11) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, imag(freq[i]), imag(x[i]))
			t.Error("FAIL")
		}
	}
}

func TestInverseDFT(t *testing.T) {
	r := NewRand64(1023)
	time := make([]complex128, 1024)
	for k := 0; k < 10; k++ {
		for i := range time {
			time[i] = r.Complex128()
		}
		freq := DFT(time)
		x := InverseDFT(freq)
		for i := range time {
			if !equal(real(x[i]), real(time[i]), 1e-9) {
				t.Logf("for i %d expected %.20f, actual %.20f\n", i, real(time[i]), real(x[i]))
				t.Error("FAIL")
			}
			if !equal(imag(x[i]), imag(time[i]), 1e-9) {
				t.Logf("for i %d expected %.20f, actual %.20f\n", i, imag(time[i]), imag(x[i]))
				t.Error("FAIL")
			}
		}
	}
}

//-----------------------------------------------------------------------------

func TestFFT(t *testing.T) {
	r := NewRand64(1023)
	in := make([]complex128, 1024)
	for k := 0; k < 10; k++ {
		for i := range in {
			in[i] = r.Complex128()
		}
		out0 := DFT(in)
		out1 := FFT(in)
		for i := range in {
			if !equal(real(out0[i]), real(out1[i]), 1e-9) {
				t.Logf("for i %d expected %.20f, actual %.20f\n", i, real(out0[i]), real(out1[i]))
				t.Error("FAIL")
			}
			if !equal(imag(out0[i]), imag(out1[i]), 1e-9) {
				t.Logf("for i %d expected %.20f, actual %.20f\n", i, imag(out0[i]), imag(out1[i]))
				t.Error("FAIL")
			}
		}
	}
}

func TestFFTx(t *testing.T) {
	r := NewRand64(1023)
	in := make([]complex128, 8)
	for k := 0; k < 10; k++ {
		for i := range in {
			in[i] = r.Complex128()
		}
		out0 := DFT(in)
		out1 := FFTx(in)
		for i := range in {
			if !equal(real(out0[i]), real(out1[i]), 1e-9) {
				t.Logf("for i %d expected %.20f, actual %.20f\n", i, real(out0[i]), real(out1[i]))
				t.Error("FAIL")
			}
			if !equal(imag(out0[i]), imag(out1[i]), 1e-9) {
				t.Logf("for i %d expected %.20f, actual %.20f\n", i, imag(out0[i]), imag(out1[i]))
				t.Error("FAIL")
			}
		}
	}
}

func TestInverseFFT(t *testing.T) {
	r := NewRand64(1023)
	in := make([]complex128, 1024)
	for k := 0; k < 10; k++ {
		for i := range in {
			in[i] = r.Complex128()
		}
		out0 := InverseDFT(in)
		out1 := InverseFFT(in)
		for i := range in {
			if !equal(real(out0[i]), real(out1[i]), 1e-9) {
				t.Logf("for i %d expected %.20f, actual %.20f\n", i, real(out0[i]), real(out1[i]))
				t.Error("FAIL")
			}
			if !equal(imag(out0[i]), imag(out1[i]), 1e-9) {
				t.Logf("for i %d expected %.20f, actual %.20f\n", i, imag(out0[i]), imag(out1[i]))
				t.Error("FAIL")
			}
		}
	}
}

//-----------------------------------------------------------------------------
// DFT benchmarking

func benchmarkDFT(n int, b *testing.B) {
	r := NewRand64(17)
	in := make([]complex128, n)
	for i := range in {
		in[i] = r.Complex128()
	}
	for n := 0; n < b.N; n++ {
		DFT(in)
	}
}

func BenchmarkDFT32(b *testing.B)   { benchmarkDFT(32, b) }
func BenchmarkDFT64(b *testing.B)   { benchmarkDFT(64, b) }
func BenchmarkDFT256(b *testing.B)  { benchmarkDFT(256, b) }
func BenchmarkDFT1024(b *testing.B) { benchmarkDFT(1024, b) }

//-----------------------------------------------------------------------------
// FFT benchmarking

func benchmarkFFT(n int, b *testing.B) {
	r := NewRand64(17)
	in := make([]complex128, n)
	for i := range in {
		in[i] = r.Complex128()
	}
	for n := 0; n < b.N; n++ {
		FFT(in)
	}
}

func BenchmarkFFT32(b *testing.B)   { benchmarkFFT(32, b) }
func BenchmarkFFT64(b *testing.B)   { benchmarkFFT(64, b) }
func BenchmarkFFT256(b *testing.B)  { benchmarkFFT(256, b) }
func BenchmarkFFT1024(b *testing.B) { benchmarkFFT(1024, b) }

//-----------------------------------------------------------------------------
