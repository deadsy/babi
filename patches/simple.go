//-----------------------------------------------------------------------------
/*

Simple Patch: an ADSR envelope on a sine wave.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/osc"
)

//-----------------------------------------------------------------------------

type simplePatch struct {
	adsr core.Module // adsr envelope
	sine core.Module // sine oscillator
}

func NewSimple() core.Module {
	log.Info.Printf("")
	return &simplePatch{
		adsr: env.NewADSR(),
		sine: osc.NewSine(),
	}
}

// Stop and performs any cleanup of a module.
func (m *simplePatch) Stop() {
	log.Info.Printf("")
	m.adsr.Stop()
	m.sine.Stop()
}

//-----------------------------------------------------------------------------
// Ports

var simplePorts = []core.PortInfo{
	{"out", "output", core.PortType_AudioBuffer, core.PortDirn_Out},
	{"f", "frequency (Hz)", core.PortType_EventFloat32, core.PortDirn_In},
}

// Ports returns the module port information.
func (m *simplePatch) Ports() []core.PortInfo {
	return simplePorts
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *simplePatch) Event(e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *simplePatch) Process(buf ...*core.Buf) {
	out := buf[0]
	m.sine.Process(out)
	var env core.Buf
	m.adsr.Process(&env)
	out.Mul(&env)
}

// Active return true if the module has non-zero output.
func (m *simplePatch) Active() bool {
	return m.adsr.Active()
}

//-----------------------------------------------------------------------------
