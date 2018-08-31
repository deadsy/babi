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
	cc   uint8       // MIDI control change number
	dst  core.Module // destination module
	gate uint        // port ID for gate of the destination module
}

// NewNote returns a MIDI note on/off gate control module.
func NewNote(ch, cc uint8, dst core.Module) core.Module {
	log.Info.Printf("")
	return &noteModule{
		ch:   ch,
		cc:   cc,
		dst:  dst,
		gate: dst.Info().GetPortID("gate"),
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
			vel := core.MIDI_Map(me.GetVelocity(), 0, 1)
			m.dst.Event(core.NewEventFloat(m.gate, vel))
		case core.EventMIDI_NoteOff:
			m.dst.Event(core.NewEventFloat(m.gate, 0))
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
