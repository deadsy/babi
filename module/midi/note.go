//-----------------------------------------------------------------------------
/*

MIDI Note Trigger Module

Generate a gate event from the MIDI note on/off events of a designated note.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *noteMidi) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "noteMidi",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortTypeMIDI, noteMidiIn},
		},
		Out: nil,
	}
}

//-----------------------------------------------------------------------------

type noteMidi struct {
	synth *core.Synth // top-level synth
	ch    uint8       // MIDI channel
	note  uint8       // MIDI note number
	dst   core.Module // destination module
	name  string      // gate port name
}

// NewNote returns a MIDI note on/off gate control module.
func NewNote(s *core.Synth, ch, note uint8, dst core.Module, name string) core.Module {
	mi := dst.Info()
	log.Info.Printf("midi ch %d note %d controlling %s.%s port", ch, note, mi.Name, name)
	return &noteMidi{
		synth: s,
		ch:    ch,
		note:  note,
		dst:   dst,
		name:  name,
	}
}

// Return the child modules.
func (m *noteMidi) Child() []core.Module {
	return nil
}

// Stop and cleanup the module.
func (m *noteMidi) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func noteMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*noteMidi)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		case core.EventMIDINoteOn:
			if me.GetNote() == m.note {
				vel := core.MIDIMap(me.GetVelocity(), 0, 1)
				core.SendEventFloat(m.dst, m.name, vel)
			}
		case core.EventMIDINoteOff:
			if me.GetNote() == m.note {
				core.SendEventFloat(m.dst, m.name, 0)
			}
		default:
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *noteMidi) Process(buf ...*core.Buf) {
	// do nothing
}

// Active return true if the module has non-zero output.
func (m *noteMidi) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
