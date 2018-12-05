//-----------------------------------------------------------------------------
/*

Low Frequency Oscillator

*/
//-----------------------------------------------------------------------------

package osc

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var lfoOscInfo = core.ModuleInfo{
	Name: "lfoOsc",
	In: []core.PortInfo{
		{"rate", "rate (Hz)", core.PortTypeFloat, lfoOscRate},
		{"depth", "depth (>= 0)", core.PortTypeFloat, lfoOscDepth},
		{"shape", "wave shape (0..5)", core.PortTypeInt, lfoOscShape},
		{"sync", "reset the lfo phase", core.PortTypeBool, lfoOscSync},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *lfoOsc) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type lfoWaveShape int

// LFO waveforms.
const (
	LfoTriangle      lfoWaveShape = 0
	LfoSawDown                    = 1
	LfoSawUp                      = 2
	LfoSquare                     = 3
	LfoSine                       = 4
	LfoSampleAndHold              = 5
)

type lfoOsc struct {
	info      core.ModuleInfo // module info
	shape     lfoWaveShape    // wave shape
	depth     float32         // wave amplitude
	x         uint32          // current x-value
	xstep     uint32          // current x-step
	randState uint32          // random state for s&h
}

// NewLFO returns a low frequency oscillator module.
func NewLFO(s *core.Synth) core.Module {
	log.Info.Printf("")
	m := &lfoOsc{
		info: lfoOscInfo,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *lfoOsc) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *lfoOsc) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func lfoOscRate(cm core.Module, e *core.Event) {
	m := cm.(*lfoOsc)
	rate := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set rate %f Hz", rate)
	m.xstep = uint32(rate * core.FrequencyScale)
}

func lfoOscShape(cm core.Module, e *core.Event) {
	m := cm.(*lfoOsc)
	shape := core.ClampInt(e.GetEventInt().Val, 0, 5)
	log.Info.Printf("set wave shape %d", shape)
	m.shape = lfoWaveShape(shape)
}

func lfoOscDepth(cm core.Module, e *core.Event) {
	m := cm.(*lfoOsc)
	depth := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set depth %f", depth)
	m.depth = depth
}

func lfoOscSync(cm core.Module, e *core.Event) {
	m := cm.(*lfoOsc)
	if e.GetEventBool().Val {
		log.Info.Printf("lfo sync")
		m.x = 0
	}
}

//-----------------------------------------------------------------------------

// Each waveform ranges from -1.0 to 1.0
// Each waveform is 0 at m.x == 0
func (m *lfoOsc) sample() float32 {
	// calculate samples as q8.24
	var sample int32
	switch m.shape {
	case LfoTriangle:
		x := m.x + (1 << 30)
		sample = int32(x >> 6)
		sample ^= -int32(x >> 31)
		sample &= (1 << 25) - 1
		sample -= (1 << 24)
	case LfoSawDown:
		sample = -int32(m.x) >> 7
	case LfoSawUp:
		sample = int32(m.x) >> 7
	case LfoSquare:
		sample = int32(m.x & (1 << 31))
		sample = (sample >> 6) | (1 << 24)
	case LfoSine:
		x := m.x - (1 << 30)
		return core.CosLookup(x)
	case LfoSampleAndHold:
		if m.x < m.xstep {
			// 0..253, cycle length = 128, 64 values with bit 7 = 1
			m.randState = ((m.randState * 179) + 17) & 0xff
		}
		sample = int32(m.randState<<24) >> 7
	}
	// convert q8.24 to float
	return float32(sample) / float32(1<<24)
}

// Process runs the module DSP.
func (m *lfoOsc) Process(buf ...*core.Buf) {
	out := buf[0]
	for i := 0; i < len(out); i++ {
		m.x += m.xstep
		out[i] = m.depth * m.sample()
	}
}

// Active returns true if the module has non-zero output.
func (m *lfoOsc) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
