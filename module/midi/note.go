//-----------------------------------------------------------------------------
/*

MIDI Note Trigger Module

Generate a gate event from the MIDI note on/off events of a designated note.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *noteModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "midi_note",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortType_EventMIDI, 0},
		},
		Out: nil,
	}
}

//-----------------------------------------------------------------------------

type noteModule struct {
	ch   uint8       // MIDI channel
	note uint8       // MIDI note number
	dst  core.Module // destination module
	gate core.PortId // gate port ID
}

// NewNote returns a MIDI note on/off gate control module.
func NewNote(ch, note uint8, dst core.Module, name string) core.Module {
	mi := dst.Info()
	log.Info.Printf("midi ch %d note %d controlling %s.%s port", ch, note, mi.Name, name)
	return &noteModule{
		ch:   ch,
		note: note,
		dst:  dst,
		gate: mi.GetPortId(name),
	}
}

// Stop and cleanup the module.
func (m *noteModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------

// Event processes a module event.
func (m *noteModule) Event(e *core.Event) {
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		case core.EventMIDI_NoteOn:
			if me.GetNote() == m.note {
				vel := core.MIDI_Map(me.GetVelocity(), 0, 1)
				core.SendEventFloatID(m.dst, m.gate, vel)
			}
		case core.EventMIDI_NoteOff:
			if me.GetNote() == m.note {
				core.SendEventFloatID(m.dst, m.gate, 0)
			}
		default:
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *noteModule) Process(buf ...*core.Buf) {
	// do nothing
}

// Active return true if the module has non-zero output.
func (m *noteModule) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
