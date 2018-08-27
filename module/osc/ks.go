//-----------------------------------------------------------------------------
/*

Karplus Strong Oscillator Module

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

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

const KS_DELAY_BITS = 6
const KS_DELAY_SIZE = 1 << KS_DELAY_BITS

// frequency to x scaling (xrange/fs)
const KS_FSCALE = (1 << 32) / core.AUDIO_FS

const KS_DELAY_MASK = KS_DELAY_SIZE - 1
const KS_FRAC_BITS = 32 - KS_DELAY_BITS
const KS_FRAC_MASK = (1 << KS_FRAC_BITS) - 1
const KS_FRAC_SCALE = 1 / (1 << KS_FRAC_BITS)

type ksModule struct {
	rand  *core.Rand
	delay [KS_DELAY_SIZE]float32 // delay line
	k     float32                // attenuation and averaging constant 0 to 0.5
	freq  float32                // base frequency
	x     uint32                 // phase position
	xstep uint32                 // phase step per sample
}

// NewKarplusStrong returns a Karplus Strong oscillator module.
func NewKarplusStrong() core.Module {
	log.Info.Printf("")
	return &ksModule{
		rand: core.NewRand(0),
	}
}

// Stop performs any cleanup of a module.
func (m *ksModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Ports

var ksPorts = []core.PortInfo{
	{"out", "output", core.PortType_Buf, core.PortDirn_Out, nil},
	{"gate", "oscillator gate, attack(>0) or mute(=0)", core.PortType_Ctrl, core.PortDirn_In, nil},
	{"f", "frequency (Hz)", core.PortType_Ctrl, core.PortDirn_In, nil},
	{"a", "attenuate (0..1)", core.PortType_Ctrl, core.PortDirn_In, nil},
}

// Ports returns the module port information.
func (m *ksModule) Ports() []core.PortInfo {
	return ksPorts
}

//-----------------------------------------------------------------------------
// Events

func (m *ksModule) event(etype oscEvent, val float32) {
	switch etype {
	case frequency: // set the oscillator frequency
		m.freq = val
		m.xstep = uint32(m.freq * KS_FSCALE)
	case attenuation: // set the attenuation
		m.k = 0.5 * core.Clamp(val, 0, 1)
	case gate: // attack(>0) or mute(=0)
		if val > 0 {
			// Initialise the delay buffer with random samples between -1 and 1.
			// The values should sum to zero so that multiple rounds of filtering
			// will make all values fall to zero.
			var sum float32
			for i := 0; i < KS_DELAY_SIZE-1; i++ {
				val := m.rand.Float()
				x := sum + val
				if x > 1 || x < -1 {
					val = -val
				}
				sum += val
				m.delay[i] = val
			}
			m.delay[KS_DELAY_SIZE-1] = -sum
		} else {
			for i := 0; i < KS_DELAY_SIZE; i++ {
				m.delay[i] = 0
			}
		}
	default:
		panic(fmt.Sprintf("unhandled event type %d", etype))
	}
}

// Event processes a module event.
func (m *ksModule) Event(e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ksModule) Process(buf []*core.Buf) {
	out := buf[0]
	for i := 0; i < len(out); i++ {
		x0 := m.x >> KS_FRAC_BITS
		x1 := (x0 + 1) & KS_DELAY_MASK
		y0 := m.delay[x0]
		y1 := m.delay[x1]
		// interpolate
		out[i] = y0 + (y1-y0)*KS_FRAC_SCALE*float32(m.x&KS_FRAC_MASK)
		// step the x position
		m.x += m.xstep
		// filter - once we have moved beyond the delay line index we
		// will average it's amplitude with the next value.
		if x0 != (m.x >> KS_FRAC_BITS) {
			m.delay[x0] = m.k * (y0 + y1)
		}
	}
}

// Active return true if the module has non-zero output.
func (m *ksModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
