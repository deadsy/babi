//-----------------------------------------------------------------------------
/*

Sequencer Testing Patch

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/seq"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *seqTestModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "seqtest_patch",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortType_EventMIDI, 0},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortType_AudioBuffer, 0},
			{"out_right", "right channel output", core.PortType_AudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type seqTestModule struct {
	synth *core.Synth // top-level synth
	seq   core.Module // sequencer
}

// NewSequencerTest returns a seqeuncer test patch.
func NewSequencerTest(s *core.Synth, prog []seq.Op) core.Module {
	log.Info.Printf("")

	sx := seq.NewSequencer(s, prog)

	// defaults
	core.SendEventFloatName(sx, "bpm", 120.0)
	core.SendEventIntName(sx, "ctrl", seq.CtrlStart)

	return &seqTestModule{
		synth: s,
		seq:   sx,
	}
}

// Child returns the child modules of this module.
func (m *seqTestModule) Child() []core.Module {
	return []core.Module{m.seq}
}

// Stop performs any cleanup of a module.
func (m *seqTestModule) Stop() {
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *seqTestModule) Event(e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *seqTestModule) Process(buf ...*core.Buf) {
	m.seq.Process(nil)
}

// Active returns true if the module has non-zero output.
func (m *seqTestModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
