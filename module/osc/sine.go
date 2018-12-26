//-----------------------------------------------------------------------------
/*

Sine Oscillator Module

*/
//-----------------------------------------------------------------------------

package osc

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var sineOscInfo = core.ModuleInfo{
	Name: "sineOsc",
	In: []core.PortInfo{
		{"frequency", "frequency (Hz)", core.PortTypeFloat, sinePortFrequency},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *sineOsc) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type sineOsc struct {
	info  core.ModuleInfo // module info
	freq  float32         // base frequency
	x     uint32          // current x-value
	xstep uint32          // current x-step
}

// NewSine returns an sine oscillator module.
func NewSine(s *core.Synth) core.Module {
	log.Info.Printf("")
	m := &sineOsc{
		info: sineOscInfo,
	}
	return s.Register(m)
}

// Return the child modules.
func (m *sineOsc) Child() []core.Module {
	return nil
}

// Stop and performs any cleanup of a module.
func (m *sineOsc) Stop() {
}

//-----------------------------------------------------------------------------
// Events

func sinePortFrequency(cm core.Module, e *core.Event) {
	m := cm.(*sineOsc)
	frequency := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set frequency %f Hz", frequency)
	m.freq = frequency
	m.xstep = uint32(frequency * core.FrequencyScale)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *sineOsc) Process(buf ...*core.Buf) bool {
	out := buf[0]
	for i := 0; i < len(out); i++ {
		out[i] = core.CosLookup(m.x)
		m.x += m.xstep
	}
	return true
}

//-----------------------------------------------------------------------------
