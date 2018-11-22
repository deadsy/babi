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

// Info returns the module information.
func (m *sineModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "sine",
		In: []core.PortInfo{
			{"frequency", "frequency (Hz)", core.PortTypeFloat, sinePortFrequency},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudio, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type sineModule struct {
	synth *core.Synth // top-level synth
	freq  float32     // base frequency
	x     uint32      // current x-value
	xstep uint32      // current x-step
}

// NewSine returns an sine oscillator module.
func NewSine(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &sineModule{
		synth: s,
	}
}

// Return the child modules.
func (m *sineModule) Child() []core.Module {
	return nil
}

// Stop and performs any cleanup of a module.
func (m *sineModule) Stop() {
}

//-----------------------------------------------------------------------------
// Events

func sinePortFrequency(cm core.Module, e *core.Event) {
	m := cm.(*sineModule)
	frequency := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set frequency %f Hz", frequency)
	m.freq = frequency
	m.xstep = uint32(frequency * core.FrequencyScale)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *sineModule) Process(buf ...*core.Buf) {
	out := buf[0]
	for i := 0; i < len(out); i++ {
		out[i] = core.CosLookup(m.x)
		m.x += m.xstep
	}
}

// Active return true if the module has non-zero output.
func (m *sineModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
