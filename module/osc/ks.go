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

var ksOscInfo = core.ModuleInfo{
	Name: "ksOsc",
	In: []core.PortInfo{
		{"gate", "oscillator gate, attack(>0) or mute(=0)", core.PortTypeFloat, ksPortGate},
		{"frequency", "frequency (Hz)", core.PortTypeFloat, ksPortFrequency},
		{"attenuation", "attenuation (0..1)", core.PortTypeFloat, ksPortAttenuation},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *ksOsc) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

const ksDelayBits = 5
const ksDelaySize = 1 << ksDelayBits

const ksDelayMask = ksDelaySize - 1
const ksFracBits = 32 - ksDelayBits
const ksFracMask = (1 << ksFracBits) - 1
const ksFracScale = 1 / (1 << ksFracBits)

type ksOsc struct {
	info  core.ModuleInfo // module info
	rand  *core.Rand32
	delay [ksDelaySize]float32 // delay line
	k     float32              // attenuation and averaging constant 0 to 0.5
	freq  float32              // base frequency
	x     uint32               // phase position
	xstep uint32               // phase step per sample
}

// NewKarplusStrong returns a Karplus Strong oscillator module.
func NewKarplusStrong(s *core.Synth) core.Module {
	log.Info.Printf("new osc")
	m := &ksOsc{
		info: ksOscInfo,
		rand: core.NewRand32(0),
	}
	return s.Register(m)
}

// Return the child modules.
func (m *ksOsc) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *ksOsc) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func ksPortGate(cm core.Module, e *core.Event) {
	m := cm.(*ksOsc)
	gate := e.GetEventFloat().Val
	log.Info.Printf("gate %f", gate)
	if gate > 0 {
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
}

func ksPortAttenuation(cm core.Module, e *core.Event) {
	m := cm.(*ksOsc)
	attenuation := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set attenuation %f", attenuation)
	m.k = 0.5 * attenuation
}

func ksPortFrequency(cm core.Module, e *core.Event) {
	m := cm.(*ksOsc)
	frequency := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set frequency %f Hz", frequency)
	m.freq = frequency
	m.xstep = uint32(frequency * core.FrequencyScale)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ksOsc) Process(buf ...*core.Buf) bool {
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
	return true
}

//-----------------------------------------------------------------------------
