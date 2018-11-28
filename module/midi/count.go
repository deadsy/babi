//-----------------------------------------------------------------------------
/*

MIDI Note Counter Module

Increment a modulo counter with each note on event.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var countMidiInfo = core.ModuleInfo{
	Name: "countMidi",
	In: []core.PortInfo{
		{"midi", "midi", core.PortTypeMIDI, countMidiIn},
	},
	Out: []core.PortInfo{
		{"count", "counter", core.PortTypeInt, nil},
	},
}

// Info returns the module information.
func (m *countMidi) Info() *core.ModuleInfo {
	return &countMidiInfo
}

// ID returns the unique module identifier.
func (m *countMidi) ID() string {
	return m.id
}

//-----------------------------------------------------------------------------

type countMidi struct {
	synth *core.Synth // top-level synth
	id    string      // module identifier
	ch    uint8       // MIDI channel
	note  uint8       // MIDI note number
	k     uint        // modulo number
	count uint        // counter
}

// NewCounter returns a MIDI counter module.
func NewCounter(s *core.Synth, ch, note uint8, k uint) core.Module {
	log.Info.Printf("")
	return &countMidi{
		synth: s,
		id:    core.GenerateID(countMidiInfo.Name),
		ch:    ch,
		note:  note,
		k:     k,
	}
}

// Child returns the child modules of this module.
func (m *countMidi) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *countMidi) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func countMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*countMidi)

	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		case core.EventMIDINoteOn:
			if me.GetNote() == m.note {

				m.count = (m.count + 1) % m.k

				// TODO send ....

			}
		default:
		}
	}

}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *countMidi) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *countMidi) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
