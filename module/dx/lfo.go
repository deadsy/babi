//-----------------------------------------------------------------------------
/*

DX7 Low Frequency Oscillator

*/
//-----------------------------------------------------------------------------

package dx

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *lfoDx) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "lfoDx",
		In:   nil,
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudio, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type lfoDx struct {
	synth *core.Synth // top-level synth
}

// NewLFO returns a DX7 low frequency oscillator module.
func NewLFO(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &lfoDx{
		synth: s,
	}
}

// Child returns the child modules of this module.
func (m *lfoDx) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *lfoDx) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *lfoDx) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *lfoDx) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
