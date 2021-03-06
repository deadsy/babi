//-----------------------------------------------------------------------------
/*

MinBLEP testing

*/
//-----------------------------------------------------------------------------

package core

import (
	"math"
	"testing"
)

//-----------------------------------------------------------------------------
// Floating Point Comparisons
// See: http://floating-point-gui.de/errors/NearlyEqualsTest.java

const minNormal = 2.2250738585072014E-308 // 2**-1022

func equal(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)
	if a == 0 || b == 0 || diff < minNormal {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diff < (epsilon * minNormal)
	}
	// use relative error
	return diff/math.Min((absA+absB), math.MaxFloat64) < epsilon
}

func convertToFloat64(ibuf []uint64) (fbuf []float64) {
	fbuf = make([]float64, len(ibuf))
	for i := 0; i < len(ibuf); i++ {
		fbuf[i] = math.Float64frombits(ibuf[i])
	}
	return
}

//-----------------------------------------------------------------------------

func TestSinc(t *testing.T) {
	xi := []uint64{
		0x401b370b70000000, 0xc000e61350000000, 0x4016a5df50000000, 0x4017e00d50000000, 0x40207744e8000000, 0xc0183225e0000000, 0xc00a5d46b0000000, 0x40157557a0000000,
		0xc011c72c78000000, 0x3ff1453880000000, 0xbfdcee8880000000, 0x40049e8d60000000, 0xc005a26d80000000, 0x3fd1273600000000, 0x402216d520000000, 0x4020a5d678000000,
		0x4005b6c0a0000000, 0x4011623de0000000, 0xc01cabfac4000000, 0x40011d7200000000, 0xc023591507400000, 0xc01491add4000000, 0xc01d057f54000000, 0x4018558a30000000,
		0xc01b773650000000, 0xbfffb2a2e0000000, 0xc01d9de44c000000, 0xc01f4b9884000000, 0x4023f4fca0000000, 0xc0168a193c000000, 0x3fd08db200000000, 0x401b2104b0000000,
	}
	yi := []uint64{
		0x3f9bb33863696c8e, 0x3faaab09f2644feb, 0xbfa922fd328f6486, 0xbf75602fd8d4a8f0, 0x3f9a74ad4f70e59b, 0x3f8083ebb9dd5c54, 0xbfb3cc7e3feab1f3, 0xbfabac0f43fa8d45,
		0x3fb20e74aa03a301, 0xbfb2a2ed8410b62b, 0x3fe646f69cb12cde, 0x3fbeaf7495c04ec6, 0x3fb822b32aa56d4a, 0x3fec59c21dbf46f5, 0xbf74212736956efd, 0x3fa0a896d9fb4eff,
		0x3fb779ae6f4ea939, 0x3fb098bbef6c001d, 0xbf96e5200580b508, 0x3fb02673a37406e6, 0xbf9cc8bc0b51a306, 0xbf9b655d4b764b3d, 0xbfa0262af7e3bd62, 0x3f8bccd0629d8684,
		0x3f935777bb2dd1ce, 0xbf83838da923116a, 0xbfa5058613e743d4, 0xbf95e622bf5fd743, 0xbf61a548c4ea4598, 0xbfaa5d832ab46d5f, 0x3fec97d03b59fb74, 0x3f9e5f4a79c9fcaa,
	}
	x := convertToFloat64(xi)
	y := convertToFloat64(yi)
	for i := 0; i < len(xi); i++ {
		result := Sinc(x[i])
		if !equal(y[i], result, 1e-12) {
			t.Logf("for x %f expected %.20f, actual %.20f\n", x[i], y[i], result)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestBlackman(t *testing.T) {
	wi := []uint64{
		0xbc70000000000000, 0x3f6ebbca113e7b00, 0x3f90038676002b96, 0x3fa3267257557c80, 0x3fb24b8125b1878f, 0x3fbecb17f77fc969, 0x3fc7c345979582a8, 0x3fd1262c9b07126b,
		0x3fd76833ef642891, 0x3fde6c9a193c12d9, 0x3fe2eb583c2a7379, 0x3fe6978ef74c6408, 0x3fe9f7c0c3383945, 0x3fecc8a026e5e20a, 0x3feeceb09ecbef93, 0x3fefdd91ca280c48,
		0x3fefdd91ca280c4c, 0x3feeceb09ecbefa2, 0x3fecc8a026e5e221, 0x3fe9f7c0c3383962, 0x3fe6978ef74c6429, 0x3fe2eb583c2a739a, 0x3fde6c9a193c1318, 0x3fd76833ef6428d1,
		0x3fd1262c9b07129a, 0x3fc7c345979582fc, 0x3fbecb17f77fc9ee, 0x3fb24b8125b187f9, 0x3fa3267257557cf6, 0x3f90038676002c32, 0x3f6ebbca113e7d40, 0xbc70000000000000,
	}
	w := convertToFloat64(wi)
	n := len(w)
	y := BlackmanWindow(n)
	for i := 0; i < n; i++ {
		if !equal(y[i], w[i], 1e-13) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, w[i], y[i])
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestRealCepstrum(t *testing.T) {

	signali := []uint64{
		0x3ff7764280000000, 0x4014724bb0000000, 0xc021ec2582000000, 0xc01b601c2c000000, 0x4023ffef20000000, 0xc017a7599c000000, 0x401f324aa0000000, 0xc01df667d4000000,
		0x4023e97628000000, 0xc021d6734f000000, 0x401da4a7f0000000, 0xc0211b5a65000000, 0xc023d56297900000, 0x4020ec3a50000000, 0x3ffe0ba540000000, 0xc01991f9d8000000,
		0xc01af31118000000, 0xc001545e50000000, 0x40208564b0000000, 0x4019935b70000000, 0xc0068b7430000000, 0x3ff0cb9340000000, 0x3ff96ae780000000, 0xbfee59f9c0000000,
		0x400dfb6360000000, 0xc02003af8d000000, 0x3fe3b79400000000, 0x40149560a0000000, 0xc00f501210000000, 0x4023b06b58000000, 0x3ff8a17940000000, 0x401e3587c0000000,
	}
	realCepstrumi := []uint64{
		0x400ae7ebe4000000, 0xbfcbee3860268bb8, 0x3f67dda2ddca15c0, 0xbfb8ff728a42059d, 0xbfa1c18b4db2921e, 0x3f982f455432530c, 0xbf7bcfe1d5995ba0, 0xbfa3d572f9536913,
		0xbfa0c17a0000005f, 0x3fc8ccd5ad06a379, 0x3f905b96ff2a59ce, 0xbf95ff54caf2ea94, 0xbfb1d6b15926b8a4, 0x3fa220b5bcdbad64, 0xbfa73f1372beaaff, 0xbfc63b90c9c9167d,
		0xbfba2a4e80000000, 0xbfc63b90c9c915a1, 0xbfa73f1372bea37e, 0x3fa220b5bcdbb258, 0xbfb1d6b15926b571, 0xbf95ff54caf2da9c, 0x3f905b96ff2a6964, 0x3fc8ccd5ad06a4ef,
		0xbfa0c179fffffcc8, 0xbfa3d572f953643a, 0xbf7bcfe1d5991628, 0x3f982f4554326328, 0xbfa1c18b4db28488, 0xbfb8ff728a41fc30, 0x3f67dda2ddcbc500, 0xbfcbee3860267cc5,
	}

	signal := convertToFloat64(signali)
	realCepstrum := convertToFloat64(realCepstrumi)

	x := RealCepstrum(signal)
	n := len(signal)
	for i := 0; i < n; i++ {
		if !equal(x[i], realCepstrum[i], 1e-6) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, realCepstrum[i], x[i])
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestMinBLEP(t *testing.T) {

	minBLEPi := []uint64{
		0xbf63f31459fe7913, 0xbf504befd8694319, 0x3f85bcf96c5402ca, 0x3fa63c7c3e149ccf, 0x3fbd080f3459c989, 0x3fce5436686d1d8b, 0x3fdb0199f282b97b, 0x3fe511a3fa904c7b,
		0x3fed417151b27a8e, 0x3ff2459267092b7a, 0x3ff4b818d92e0115, 0x3ff58c06c6562e56, 0x3ff4dc055e2c8adf, 0x3ff342af9a0631f8, 0x3ff1979dcf8aad22, 0x3ff090e1643eb41e,
		0x3ff0797fe63516dc, 0x3ff121881af7ccd5, 0x3ff20b435f9236ce, 0x3ff2b752818e5225, 0x3ff2e3b26747e233, 0x3ff29e9a83e2be12, 0x3ff22b704a6543ad, 0x3ff1d1faa2bf37a3,
		0x3ff1b7464ef4e049, 0x3ff1d2f2a151afe0, 0x3ff1fec5626f21e0, 0x3ff21053cca06ae2, 0x3ff1ec6bda2ac09a, 0x3ff18db5984fecb7, 0x3ff101e4d36d92a2, 0x3ff063bd419491ad,
		0x3fefa95831b85d68, 0x3feee974538aa429, 0x3feeb007089508e4, 0x3feefa4838a783cd, 0x3fef963d4e5d4d0e, 0x3ff01ad626d765da, 0x3ff048f8d2fc048c, 0x3ff046a989c43980,
		0x3ff01ee554dbd383, 0x3fefdaf54a81fee1, 0x3fef9b5189eb8eaa, 0x3fef96f313324d71, 0x3fefc06f2d1de742, 0x3feff4c884bac2d9, 0x3ff00a445baad4b1, 0x3ff00a6d1db446db,
		0x3ff0000000000000,
	}

	minBLEP := convertToFloat64(minBLEPi)

	x := GenerateMinBLEP(6, 4)
	n := len(minBLEP)
	for i := 0; i < n; i++ {
		if !equal(x[i], minBLEP[i], 1e-6) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, minBLEP[i], x[i])
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------
/*

// Reference C++ code to generate test vectors.
// MinBLEP Generation Code
// By Daniel Werner
// This Code Is Public Domain

#include <math.h>
#include <stdio.h>
#include <stdlib.h>
#include <inttypes.h>

#define PI 3.14159265358979

// SINC Function
inline double SINC(double x) {
	double pix;

	if (x == 0.0)
		return 1.0;
	else {
		pix = PI * x;
		return sin(pix) / pix;
	}
}

// Generate Blackman Window
inline void BlackmanWindow(int n, double *w) {
	int m = n - 1;
	int i;
	double f1, f2, fm;

	fm = (double)m;
	for (i = 0; i <= m; i++) {
		f1 = (2.0 * PI * (double)i) / fm;
		f2 = 2.0 * f1;
		w[i] = 0.42 - (0.5 * cos(f1)) + (0.08 * cos(f2));
	}
}

// Discrete Fourier Transform
void DFT(int n, double *realTime, double *imagTime, double *realFreq, double *imagFreq) {
	int k, i;
	double sr, si, p;

	for (k = 0; k < n; k++) {
		realFreq[k] = 0.0;
		imagFreq[k] = 0.0;
	}

	for (k = 0; k < n; k++)
		for (i = 0; i < n; i++) {
			p = (2.0 * PI * (double)(k * i)) / n;
			sr = cos(p);
			si = -sin(p);
			realFreq[k] += (realTime[i] * sr) - (imagTime[i] * si);
			imagFreq[k] += (realTime[i] * si) + (imagTime[i] * sr);
		}
}

// Inverse Discrete Fourier Transform
void InverseDFT(int n, double *realTime, double *imagTime, double *realFreq, double *imagFreq) {
	int k, i;
	double sr, si, p;

	for (k = 0; k < n; k++) {
		realTime[k] = 0.0;
		imagTime[k] = 0.0;
	}

	for (k = 0; k < n; k++) {
		for (i = 0; i < n; i++) {
			p = (2.0 * PI * (double)(k * i)) / n;
			sr = cos(p);
			si = -sin(p);
			realTime[k] += (realFreq[i] * sr) + (imagFreq[i] * si);
			imagTime[k] += (realFreq[i] * si) - (imagFreq[i] * sr);
		}
		realTime[k] /= n;
		imagTime[k] /= n;
	}
}

// Complex Absolute Value
inline double cabs(double x, double y) {
	return sqrt((x * x) + (y * y));
}

// Complex Exponential
inline void cexp(double x, double y, double *zx, double *zy) {
	double expx;

	expx = exp(x);
	*zx = expx * cos(y);
	*zy = expx * sin(y);
}

// Compute Real Cepstrum Of Signal
void RealCepstrum(int n, double *signal, double *realCepstrum) {
	double *realTime, *imagTime, *realFreq, *imagFreq;
	int i;

	realTime = new double[n];
	imagTime = new double[n];
	realFreq = new double[n];
	imagFreq = new double[n];

	// Compose Complex FFT Input

	for (i = 0; i < n; i++) {
		realTime[i] = signal[i];
		imagTime[i] = 0.0;
	}

	// Perform DFT

	DFT(n, realTime, imagTime, realFreq, imagFreq);

	// Calculate Log Of Absolute Value

	for (i = 0; i < n; i++) {
		realFreq[i] = logf(cabs(realFreq[i], imagFreq[i]));
		imagFreq[i] = 0.0;
	}

	// Perform Inverse FFT

	InverseDFT(n, realTime, imagTime, realFreq, imagFreq);

	// Output Real Part Of FFT
	for (i = 0; i < n; i++)
		realCepstrum[i] = realTime[i];

	delete realTime;
	delete imagTime;
	delete realFreq;
	delete imagFreq;
}

// Compute Minimum Phase Reconstruction Of Signal
void MinimumPhase(int n, double *realCepstrum, double *minimumPhase) {
	int i, nd2;
	double *realTime, *imagTime, *realFreq, *imagFreq;

	nd2 = n / 2;
	realTime = new double[n];
	imagTime = new double[n];
	realFreq = new double[n];
	imagFreq = new double[n];

	if ((n % 2) == 1) {
		realTime[0] = realCepstrum[0];
		for (i = 1; i < nd2; i++)
			realTime[i] = 2.0 * realCepstrum[i];
		for (i = nd2; i < n; i++)
			realTime[i] = 0.0;
	} else {
		realTime[0] = realCepstrum[0];
		for (i = 1; i < nd2; i++)
			realTime[i] = 2.0 * realCepstrum[i];
		realTime[nd2] = realCepstrum[nd2];
		for (i = nd2 + 1; i < n; i++)
			realTime[i] = 0.0;
	}

	for (i = 0; i < n; i++)
		imagTime[i] = 0.0;

	DFT(n, realTime, imagTime, realFreq, imagFreq);

	for (i = 0; i < n; i++)
		cexp(realFreq[i], imagFreq[i], &realFreq[i], &imagFreq[i]);

	InverseDFT(n, realTime, imagTime, realFreq, imagFreq);

	for (i = 0; i < n; i++)
		minimumPhase[i] = realTime[i];

	delete realTime;
	delete imagTime;
	delete realFreq;
	delete imagFreq;
}

// Generate MinBLEP And Return It In An Array Of Floating Point Values
double *GenerateMinBLEP(int zeroCrossings, int overSampling) {
	int i, n;
	double r, a, b;
	double *buffer1, *buffer2, *minBLEP;

	n = (zeroCrossings * 2 * overSampling) + 1;

	buffer1 = new double[n];
	buffer2 = new double[n];

	// Generate Sinc

	a = (double)-zeroCrossings;
	b = (double)zeroCrossings;
	for (i = 0; i < n; i++) {
		r = ((double)i) / ((double)(n - 1));
		buffer1[i] = SINC(a + (r * (b - a)));
	}

	// Window Sinc

	BlackmanWindow(n, buffer2);
	for (i = 0; i < n; i++)
		buffer1[i] *= buffer2[i];

	// Minimum Phase Reconstruction

	RealCepstrum(n, buffer1, buffer2);
	MinimumPhase(n, buffer2, buffer1);

	// Integrate Into MinBLEP

	minBLEP = new double[n];
	a = 0.0;
	for (i = 0; i < n; i++) {
		a += buffer1[i];
		minBLEP[i] = a;
	}

	// Normalize
	a = minBLEP[n - 1];
	a = 1.0 / a;
	for (i = 0; i < n; i++)
		minBLEP[i] *= a;

	delete buffer1;
	delete buffer2;
	return minBLEP;
}

static void dump_buf(const char *str, double *buf, size_t n) {
	if (str) {
		printf("%s := []uint64{\n", str);
	}
	for (int i = 0; i < n; i++) {
		uint64_t *ptr = (uint64_t *) & buf[i];
		//printf("%.10f,", buf[i]);
		printf("0x%016lx,", *ptr);
		if ((i != 0) && (i % 8 == 7)) {
			printf("\n");
		}
	}
	printf("}\n");
}

static double randMToN(double m, double n) {
	double x = rand() / (RAND_MAX + 1.f);
	return m + x * (n - m);
}

static void rand_buf(double *buf, size_t n) {
	for (int i = 0; i < n; i++) {
		buf[i] = randMToN(-10.0, 10.0);
	}
}

#define WINDOW_SIZE 32

int main(void) {
	double w[WINDOW_SIZE];
	double realTime[WINDOW_SIZE];
	double imagTime[WINDOW_SIZE];
	double realFreq[WINDOW_SIZE];
	double imagFreq[WINDOW_SIZE];
	double signal[WINDOW_SIZE];
	double realCepstrum[WINDOW_SIZE];
	double minimumPhase[WINDOW_SIZE];
	double x[WINDOW_SIZE];
	double y[WINDOW_SIZE];

	// sinc
	rand_buf(x, WINDOW_SIZE);
	for (int i = 0; i < WINDOW_SIZE; i++) {
		y[i] = SINC(x[i]);
	}
	dump_buf("xi", x, WINDOW_SIZE);
	dump_buf("yi", y, WINDOW_SIZE);

	// blackman
	BlackmanWindow(WINDOW_SIZE, w);
	dump_buf("window", w, WINDOW_SIZE);

	// dft
	rand_buf(realTime, WINDOW_SIZE);
	rand_buf(imagTime, WINDOW_SIZE);
	DFT(WINDOW_SIZE, realTime, imagTime, realFreq, imagFreq);
	dump_buf("realTimei", realTime, WINDOW_SIZE);
	dump_buf("imagTimei", imagTime, WINDOW_SIZE);
	dump_buf("realFreqi", realFreq, WINDOW_SIZE);
	dump_buf("imagFreqi", imagFreq, WINDOW_SIZE);

	// idft
	rand_buf(realFreq, WINDOW_SIZE);
	rand_buf(imagFreq, WINDOW_SIZE);
	InverseDFT(WINDOW_SIZE, realTime, imagTime, realFreq, imagFreq);
	dump_buf("realFreqi", realFreq, WINDOW_SIZE);
	dump_buf("imagFreqi", imagFreq, WINDOW_SIZE);
	dump_buf("realTimei", realTime, WINDOW_SIZE);
	dump_buf("imagTimei", imagTime, WINDOW_SIZE);

	// real cepstrum
	rand_buf(signal, WINDOW_SIZE);
	RealCepstrum(WINDOW_SIZE, signal, realCepstrum);
	dump_buf("signal", signal, WINDOW_SIZE);
	dump_buf("realCepstrum", realCepstrum, WINDOW_SIZE);

	// minimum phase
	rand_buf(realCepstrum, WINDOW_SIZE);
	MinimumPhase(WINDOW_SIZE, realCepstrum, minimumPhase);
	dump_buf("realCepstrum", realCepstrum, WINDOW_SIZE);
	dump_buf("minimumPhase", minimumPhase, WINDOW_SIZE);

	// MinBLEP
	double *minBLEP = GenerateMinBLEP(6, 4);
	dump_buf("minBLEP", minBLEP, 49);

	return 0;
}

*/
//-----------------------------------------------------------------------------
