//-----------------------------------------------------------------------------
/*

Sequencer Testing Patch

*/
//-----------------------------------------------------------------------------

package patch

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/seq"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

/*
var metronome = []seq.Op{
	seq.OpNote(1, 69, 100, 4),
	seq.OpRest(12),
	seq.OpNote(1, 60, 100, 4),
	seq.OpRest(12),
	seq.OpNote(1, 60, 100, 4),
	seq.OpRest(12),
	seq.OpNote(1, 60, 100, 4),
	seq.OpRest(12),
	seq.OpLoop(),
}
*/

//-----------------------------------------------------------------------------

var seqPatchInfo = core.ModuleInfo{
	Name: "seqPatch",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, seqtestPortMidiIn},
	},
	Out: []core.PortInfo{
		{"out0", "left channel output", core.PortTypeAudio, nil},
		{"out1", "right channel output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *seqPatch) Info() *core.ModuleInfo {
	return &seqPatchInfo
}

// ID returns the unique module identifier.
func (m *seqPatch) ID() string {
	return m.id
}

//-----------------------------------------------------------------------------

type seqPatch struct {
	synth *core.Synth // top-level synth
	id    string      // module identifier
	seq   core.Module // sequencer
}

// NewSequencerTest returns a seqeuncer test patch.
func NewSequencerTest(s *core.Synth, prog []seq.Op) core.Module {
	log.Info.Printf("")

	sx := seq.NewSequencer(s, prog)

	// defaults
	core.SendEventFloat(sx, "bpm", 120.0)
	core.SendEventInt(sx, "ctrl", seq.CtrlStart)

	return &seqPatch{
		synth: s,
		id:    core.GenerateID(seqPatchInfo.Name),
		seq:   sx,
	}
}

// Child returns the child modules of this module.
func (m *seqPatch) Child() []core.Module {
	return []core.Module{m.seq}
}

// Stop performs any cleanup of a module.
func (m *seqPatch) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func seqtestPortMidiIn(cm core.Module, e *core.Event) {
	// nothing...
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *seqPatch) Process(buf ...*core.Buf) {
	m.seq.Process(nil)
}

// Active returns true if the module has non-zero output.
func (m *seqPatch) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
