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

const (
	goomPortNull = iota
	goomPortFrequency
	goomPortDuty
	goomPortSlope
)

// Info returns the module information.
func (m *goomModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "goom",
		In: []core.PortInfo{
			{"fm", "frequency modulation", core.PortTypeAudioBuffer, 0},
			{"frequency", "frequency (Hz)", core.PortTypeFloat, goomPortFrequency},
			{"duty", "duty cycle (0..1)", core.PortTypeFloat, goomPortDuty},
			{"slope", "slope (0..1)", core.PortTypeFloat, goomPortSlope},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

const goomFullCycle = (1 << 32) - 1
const goomHalfCycle = 1 << 31

// Limit how close the duty cycle can get to 0/100%.
const tpMin = 0.05

// Limit how fast the slope can rise.
const slopeMin = 0.1

type goomModule struct {
	synth *core.Synth // top-level synth
	freq  float32     // base frequency
	tp    uint32      // s0f0 to s1f1 transition point
	k0    float32     // scaling factor for slope 0
	k1    float32     // scaling factor for slope 1
	x     uint32      // phase position
	xstep uint32      // phase step per sample
}

// NewGoom returns a goom oscillator module.
func NewGoom(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &goomModule{
		synth: s,
	}
}

// Child returns the child modules of this module.
func (m *goomModule) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *goomModule) Stop() {
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *goomModule) Event(e *core.Event) {
	fe := e.GetEventFloat()
	if fe != nil {
		val := fe.Val
		switch fe.ID {
		case goomPortFrequency: // set the oscillator frequency
			log.Info.Printf("set frequency %f", val)
			m.freq = val
			m.xstep = uint32(val * core.FrequencyScale)
		case goomPortDuty: // set the wave duty cycle
			log.Info.Printf("set duty cycle %f", val)
			duty := core.Clamp(val, 0, 1)
			m.tp = uint32(goomFullCycle * core.Map(duty, tpMin, 0.5))
		case goomPortSlope: // set the wave slope
			log.Info.Printf("set slope %f", val)
			// Work out the portion of s0f0/s1f1 that is sloped.
			slope := core.Clamp(val, 0, 1)
			slope = core.Map(slope, slopeMin, 1)
			// scaling constant for s0, map the slope to the LUT.
			m.k0 = 1.0 / (float32(m.tp) * slope)
			// scaling constant for s1, map the slope to the LUT.
			m.k1 = 1.0 / (float32(goomFullCycle-m.tp) * slope)
		default:
			log.Info.Printf("bad port number %d", fe.ID)
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *goomModule) Process(buf ...*core.Buf) {
	fm := buf[0]
	out := buf[1]
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
			ofs = goomHalfCycle
		}
		// clamp x to 1
		if x > 1 {
			x = 1
		}
		out[i] = core.CosLookup(uint32(x*float32(goomHalfCycle)) + ofs)
		// step the phase
		if fm != nil {
			m.x += uint32((m.freq + fm[i]) * core.FrequencyScale)
		} else {
			m.x += m.xstep
		}
	}
}

// Active returns true if the module has non-zero output.
func (m *goomModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
