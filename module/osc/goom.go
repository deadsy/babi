//-----------------------------------------------------------------------------
/*

Goom Waves

A Goom Wave is a wave shape with the following segments:

1) s0: A falling (1 to -1) sine curve
2) f0: A flat piece at the bottom
3) s1: A rising (-1 to 1) sine curve
4) f1: A flat piece at the top

Shape is controlled by two parameters:
duty = split the total period between s0,f0 and s1,f1
slope = split s0f0 and s1f1 between slope and flat.

The idea for goom waves comes from: https://www.quinapalus.com/goom.html

*/
//-----------------------------------------------------------------------------

package osc

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *goomOsc) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "goom",
		In: []core.PortInfo{
			{"frequency", "frequency (Hz)", core.PortTypeFloat, goomPortFrequency},
			{"duty", "duty cycle (0..1)", core.PortTypeFloat, goomPortDuty},
			{"slope", "slope (0..1)", core.PortTypeFloat, goomPortSlope},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type goomOsc struct {
	synth *core.Synth // top-level synth
	freq  float32     // base frequency
	tp    uint32      // s0f0 to s1f1 transition point
	k0    float32     // scaling factor for slope 0
	k1    float32     // scaling factor for slope 1
	x     uint32      // phase position
	xstep uint32      // phase step per sample
}

// NewGoomOsc returns a goom oscillator module.
func NewGoomOsc(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &goomOsc{
		synth: s,
	}
}

// Child returns the child modules of this module.
func (m *goomOsc) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *goomOsc) Stop() {
}

//-----------------------------------------------------------------------------
// Events

func goomPortFrequency(cm core.Module, e *core.Event) {
	m := cm.(*goomOsc)
	frequency := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set frequency %f Hz", frequency)
	m.freq = frequency
	m.xstep = uint32(frequency * core.FrequencyScale)
}

func goomPortDuty(cm core.Module, e *core.Event) {
	m := cm.(*goomOsc)
	duty := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set duty cycle %f", duty)
	m.tp = uint32(float32(core.FullCycle) * core.Map(duty, 0.05, 0.5))
}

func goomPortSlope(cm core.Module, e *core.Event) {
	m := cm.(*goomOsc)
	slope := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set slope %f", slope)
	// Work out the portion of s0f0/s1f1 that is sloped.
	slope = core.Map(slope, 0.1, 1)
	// scaling constant for s0, map the slope to the LUT.
	m.k0 = 1.0 / (float32(m.tp) * slope)
	// scaling constant for s1, map the slope to the LUT.
	m.k1 = 1.0 / (float32(core.FullCycle-1-m.tp) * slope)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *goomOsc) Process(buf ...*core.Buf) {
	out := buf[0]
	for i := 0; i < len(out); i++ {
		var ofs uint32
		var x float32
		// what portion of the goom wave are we in?
		if m.x < m.tp {
			// we are in the s0/f0 portion
			x = float32(m.x) * m.k0
		} else {
			// we are in the s1/f1 portion
			x = float32(m.x-m.tp) * m.k1
			ofs = core.HalfCycle
		}
		// clamp x to 1
		if x > 1 {
			x = 1
		}
		out[i] = core.CosLookup(uint32(x*float32(core.HalfCycle)) + ofs)
		// step the phase
		m.x += m.xstep
	}
}

// Active returns true if the module has non-zero output.
func (m *goomOsc) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
