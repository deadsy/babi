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
