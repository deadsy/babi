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

// Info returns the module information.
func (m *seqtestModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "seqtest",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortTypeMIDI, seqtestPortMidiIn},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortTypeAudioBuffer, nil},
			{"out_right", "right channel output", core.PortTypeAudioBuffer, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type seqtestModule struct {
	synth *core.Synth // top-level synth
	seq   core.Module // sequencer
}

// NewSequencerTest returns a seqeuncer test patch.
func NewSequencerTest(s *core.Synth, prog []seq.Op) core.Module {
	log.Info.Printf("")

	sx := seq.NewSequencer(s, prog)

	// defaults
	core.SendEventFloat(sx, "bpm", 120.0)
	core.SendEventInt(sx, "ctrl", seq.CtrlStart)

	return &seqtestModule{
		synth: s,
		seq:   sx,
	}
}

// Child returns the child modules of this module.
func (m *seqtestModule) Child() []core.Module {
	return []core.Module{m.seq}
}

// Stop performs any cleanup of a module.
func (m *seqtestModule) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func seqtestPortMidiIn(cm core.Module, e *core.Event) {
	// nothing...
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *seqtestModule) Process(buf ...*core.Buf) {
	m.seq.Process(nil)
}

// Active returns true if the module has non-zero output.
func (m *seqtestModule) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
