//-----------------------------------------------------------------------------
/*

Sine Oscillator Module

*/
//-----------------------------------------------------------------------------

package osc

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

const (
	sine_port_null = iota
	sine_port_frequency
)

// Info returns the module information.
func (m *sineModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "sine",
		In: []core.PortInfo{
			{"frequency", "frequency (Hz)", core.PortType_EventFloat, sine_port_frequency},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortType_AudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

// frequency to x scaling (xrange/fs)
const sine_freq_scale = (1 << 32) / core.AUDIO_FS

type sineModule struct {
	freq  float32 // base frequency
	x     uint32  // current x-value
	xstep uint32  // current x-step
}

// NewSine returns an sine oscillator module.
func NewSine() core.Module {
	log.Info.Printf("")
	return &sineModule{}
}

// Stop and performs any cleanup of a module.
func (m *sineModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *sineModule) Event(e *core.Event) {
	fe := e.GetEventFloat()
	if fe != nil {
		val := fe.Val
		switch fe.Id {
		case sine_port_frequency: // set the oscillator frequency
			log.Info.Printf("set frequency %f", val)
			m.freq = val
			m.xstep = uint32(val * sine_freq_scale)
		default:
			log.Info.Printf("bad port number %d", fe.Id)
		}
	}
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
