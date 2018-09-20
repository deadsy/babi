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
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const (
	ksPortNull = iota
	ksPortGate
	ksPortAttenuation
	ksPortFrequency
)

// Info returns the module information.
func (m *ksModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "karplus_strong",
		In: []core.PortInfo{
			{"gate", "oscillator gate, attack(>0) or mute(=0)", core.PortTypeFloat, ksPortGate},
			{"frequency", "frequency (Hz)", core.PortTypeFloat, ksPortFrequency},
			{"attenuation", "attenuation (0..1)", core.PortTypeFloat, ksPortAttenuation},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

const ksDelayBits = 6
const ksDelaySize = 1 << ksDelayBits

const ksDelayMask = ksDelaySize - 1
const ksFracBits = 32 - ksDelayBits
const ksFracMask = (1 << ksFracBits) - 1
const ksFracScale = 1 / (1 << ksFracBits)

type ksModule struct {
	synth *core.Synth // top-level synth
	rand  *core.Rand32
	delay [ksDelaySize]float32 // delay line
	k     float32              // attenuation and averaging constant 0 to 0.5
	freq  float32              // base frequency
	x     uint32               // phase position
	xstep uint32               // phase step per sample
}

// NewKarplusStrong returns a Karplus Strong oscillator module.
func NewKarplusStrong(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &ksModule{
		synth: s,
		rand:  core.NewRand32(0),
	}
}

// Return the child modules.
func (m *ksModule) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *ksModule) Stop() {
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *ksModule) Event(e *core.Event) {
	fe := e.GetEventFloat()
	if fe != nil {
		val := fe.Val
		switch fe.ID {
		case ksPortGate: // attack(>0) or mute(=0)
			log.Info.Printf("gate %f", val)
			if val > 0 {
				// Initialise the delay buffer with random samples between -1 and 1.
				// The values should sum to zero so that multiple rounds of filtering
				// will make all values fall to zero.
				var sum float32
				for i := 0; i < ksDelaySize-1; i++ {
					val := m.rand.Float32()
					x := sum + val
					if x > 1 || x < -1 {
						val = -val
					}
					sum += val
					m.delay[i] = val
				}
				m.delay[ksDelaySize-1] = -sum
			} else {
				for i := 0; i < ksDelaySize; i++ {
					m.delay[i] = 0
				}
			}
		case ksPortAttenuation: // set the attenuation
			log.Info.Printf("set attenuation %f", val)
			m.k = 0.5 * core.Clamp(val, 0, 1)
		case ksPortFrequency: // set the oscillator frequency
			log.Info.Printf("set frequency %f", val)
			m.freq = val
			m.xstep = uint32(val * core.FrequencyScale)
		default:
			log.Info.Printf("bad port number %d", fe.ID)
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ksModule) Process(buf ...*core.Buf) {
	out := buf[0]
	for i := 0; i < len(out); i++ {
		x0 := m.x >> ksFracBits
		x1 := (x0 + 1) & ksDelayMask
		y0 := m.delay[x0]
		y1 := m.delay[x1]
		// interpolate
		out[i] = y0 + (y1-y0)*ksFracScale*float32(m.x&ksFracMask)
		// step the x position
		m.x += m.xstep
		// filter - once we have moved beyond the delay line index we
		// will average it's amplitude with the next value.
		if x0 != (m.x >> ksFracBits) {
			m.delay[x0] = m.k * (y0 + y1)
		}
	}
}

// Active return true if the module has non-zero output.
func (m *ksModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
