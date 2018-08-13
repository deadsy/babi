//-----------------------------------------------------------------------------
/*

Karplus Strong Plucked String Modelling

KS generally has a delay line buffer size that determines the fundamental frequency
of the sound. That has some practical problems. The delay line buffer is too
large for low frequencies and it makes it hard to provide fine resolution
control over the frequency. This implementation uses a fixed buffer size and
steps through it with a 32 bit phase value. The step size determines the
frequency of the sound. When the step position falls between samples we do
linear interpolation to get the output value. When we move beyond a sample
we do the low pass filtering on it (in this case simple averaging).

*/
//-----------------------------------------------------------------------------

package osc

import "github.com/deadsy/babi/core"

//-----------------------------------------------------------------------------

const KS_DELAY_BITS = 6
const KS_DELAY_SIZE = 1 << KS_DELAY_BITS

// frequency to x scaling (xrange/fs)
const KS_FSCALE = (1 << 32) / core.AUDIO_FS

const KS_DELAY_MASK = KS_DELAY_SIZE - 1
const KS_FRAC_BITS = 32 - KS_DELAY_BITS
const KS_FRAC_MASK = (1 << KS_FRAC_BITS) - 1
const KS_FRAC_SCALE = 1 / (1 << KS_FRAC_BITS)

//-----------------------------------------------------------------------------

type KarplusStrong struct {
	rand  *core.Rand
	delay [KS_DELAY_SIZE]float32 // delay line
	k     float32                // attenuation and averaging constant 0 to 0.5
	freq  float32                // base frequency
	x     uint32                 // phase position
	xstep uint32                 // phase step per sample
}

func NewKarplusStrong() *KarplusStrong {
	return &KarplusStrong{
		rand: core.NewRand(0),
	}
}

func (ks *KarplusStrong) SetAttenuate(attenuate float32) {
	ks.k = 0.5 * core.Clamp(attenuate, 0, 1)
}

func (ks *KarplusStrong) SetFrequency(freq float32) {
	ks.freq = freq
	ks.xstep = uint32(ks.freq * KS_FSCALE)
}

func (ks *KarplusStrong) Pluck() {
	// Initialise the delay buffer with random samples between -1 and 1.
	// The values should sum to zero so that multiple rounds of filtering
	// will make all values fall to zero.
	var sum float32
	for i := 0; i < KS_DELAY_SIZE-1; i++ {
		val := ks.rand.Float()
		x := sum + val
		if x > 1 || x < -1 {
			val = -val
		}
		sum += val
		ks.delay[i] = val
	}
	ks.delay[KS_DELAY_SIZE-1] = -sum
}

func (ks *KarplusStrong) Process(out *core.Buf) {
	for i := 0; i < len(out); i++ {
		x0 := ks.x >> KS_FRAC_BITS
		x1 := (x0 + 1) & KS_DELAY_MASK
		y0 := ks.delay[x0]
		y1 := ks.delay[x1]
		// interpolate
		out[i] = y0 + (y1-y0)*KS_FRAC_SCALE*float32(ks.x&KS_FRAC_MASK)
		// step the x position
		ks.x += ks.xstep
		// filter - once we have moved beyond the delay line index we
		// will average it's amplitude with the next value.
		if x0 != (ks.x >> KS_FRAC_BITS) {
			ks.delay[x0] = ks.k * (y0 + y1)
		}
	}
}

//-----------------------------------------------------------------------------
