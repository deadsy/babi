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
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var goomOscInfo = core.ModuleInfo{
	Name: "goomOsc",
	In: []core.PortInfo{
		{"frequency", "frequency (Hz)", core.PortTypeFloat, goomOscFrequency},
		{"duty", "duty cycle (0..1)", core.PortTypeFloat, goomOscDuty},
		{"slope", "slope (0..1)", core.PortTypeFloat, goomOscSlope},
		{"mode", "oscillator mode", core.PortTypeInt, goomOscMode},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *goomOsc) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type goomMode int

// goom oscillator modes
const (
	GoomModeBasic goomMode = 0 // no feedback, no modulation
	GoomModeFM             = 1 // frequency modulation input
	GoomModePM             = 2 // phase modulation input
)

type goomOsc struct {
	info  core.ModuleInfo // module info
	mode  goomMode        // oscillator mode
	freq  float32         // base frequency
	duty  float32         // wave duty cycle
	slope float32         // wave slope
	tp    uint32          // s0f0 to s1f1 transition point
	k0    float32         // scaling factor for slope 0
	k1    float32         // scaling factor for slope 1
	x     uint32          // phase position
	xstep uint32          // phase step per sample
}

// NewGoom returns a goom oscillator module.
func NewGoom(s *core.Synth) core.Module {
	log.Info.Printf("")
	m := &goomOsc{
		info: goomOscInfo,
	}
	return s.Register(m)
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

func (m *goomOsc) setFrequency(frequency float32) {
	m.freq = frequency
	m.xstep = uint32(frequency * core.FrequencyScale)
}

func (m *goomOsc) setShape(duty, slope float32) {
	// update duty cycle
	m.duty = duty
	m.tp = uint32(float32(core.FullCycle) * core.MapLin(duty, 0.05, 0.5))
	// update the slope
	m.slope = slope
	// Work out the portion of s0f0/s1f1 that is sloped.
	slope = core.MapLin(slope, 0.1, 1)
	// scaling constant for s0, map the slope to the LUT.
	m.k0 = 1.0 / (float32(m.tp) * slope)
	// scaling constant for s1, map the slope to the LUT.
	m.k1 = 1.0 / (float32(core.FullCycle-1-m.tp) * slope)
}

func goomOscFrequency(cm core.Module, e *core.Event) {
	m := cm.(*goomOsc)
	frequency := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set frequency %f Hz", frequency)
	m.setFrequency(frequency)
}

func goomOscMode(cm core.Module, e *core.Event) {
	m := cm.(*goomOsc)
	mode := e.GetEventInt().Val
	log.Info.Printf("set mode %d", mode)
	m.mode = goomMode(mode)
}

func goomOscDuty(cm core.Module, e *core.Event) {
	m := cm.(*goomOsc)
	duty := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set duty cycle %f", duty)
	m.setShape(duty, m.slope)
}

func goomOscSlope(cm core.Module, e *core.Event) {
	m := cm.(*goomOsc)
	slope := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set slope %f", slope)
	m.setShape(m.duty, slope)
}

//-----------------------------------------------------------------------------

// sample returns a single sample for the current phase value.
func (m *goomOsc) sample() float32 {
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
	return core.CosLookup(uint32(x*float32(core.HalfCycle)) + ofs)
}

// Process runs the module DSP.
func (m *goomOsc) Process(buf ...*core.Buf) {

	switch m.mode {
	case GoomModeBasic: // no feedback, no modulation
		out := buf[0]
		for i := 0; i < len(out); i++ {
			out[i] = m.sample()
			// step the phase
			m.x += m.xstep
		}
	case GoomModeFM: // frequency modulation input
		fm := buf[0]
		out := buf[1]
		for i := 0; i < len(out); i++ {
			out[i] = m.sample()
			// step the phase
			m.x += uint32((m.freq + fm[i]) * core.FrequencyScale)
		}
	case GoomModePM: // phase modulation input
		pm := buf[0]
		out := buf[1]
		for i := 0; i < len(out); i++ {
			out[i] = m.sample()
			// step the phase
			m.x += uint32(float32(m.xstep) + (pm[i] * core.PhaseScale))
		}
	default:
		panic(fmt.Sprintf("bad mode %d", m.mode))
	}
}

// Active returns true if the module has non-zero output.
func (m *goomOsc) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
