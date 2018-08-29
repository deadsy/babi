//-----------------------------------------------------------------------------
/*

Sine Oscillator Module

*/
//-----------------------------------------------------------------------------

package osc

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *sineModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		In: []core.PortInfo{
			{"f", "frequency (Hz)", core.PortType_EventFloat},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortType_AudioBuffer},
		},
	}
}

//-----------------------------------------------------------------------------

// frequency to x scaling (xrange/fs)
const FREQ_SCALE = (1 << 32) / core.AUDIO_FS

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

func (m *sineModule) event(etype oscEvent, val float32) {
	switch etype {
	case frequency: // set the oscillator frequency
		m.freq = val
		m.xstep = uint32(m.freq * FREQ_SCALE)
	default:
		panic(fmt.Sprintf("unhandled event type %d", etype))
	}
}

// Event processes a module event.
func (m *sineModule) Event(e *core.Event) {
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
