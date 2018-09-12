//-----------------------------------------------------------------------------
/*

Sequencer

*/
//-----------------------------------------------------------------------------

package seq

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *seqModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "seq",
		In:   nil,
		Out:  nil,
	}
}

//-----------------------------------------------------------------------------

type seqModule struct {
	synth *core.Synth // top-level synth
}

// NewSeq returns a sequencer.
func NewSeq(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &seqModule{}
}

// Return the child modules.
func (m *seqModule) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *seqModule) Stop() {
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *seqModule) Event(e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *seqModule) Process(buf ...*core.Buf) {
}

// Active return true if the module has non-zero output.
func (m *seqModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
