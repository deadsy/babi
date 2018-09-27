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
func (m *noteModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "midi_note",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortTypeMIDI, notePortMidiIn},
		},
		Out: nil,
	}
}

//-----------------------------------------------------------------------------

type noteModule struct {
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
	return &noteModule{
		synth: s,
		ch:    ch,
		note:  note,
		dst:   dst,
		name:  name,
	}
}

// Return the child modules.
func (m *noteModule) Child() []core.Module {
	return nil
}

// Stop and cleanup the module.
func (m *noteModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Port Events

func notePortMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*noteModule)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		case core.EventMIDINoteOn:
			if me.GetNote() == m.note {
				vel := core.MIDIMap(me.GetVelocity(), 0, 1)
				core.SendEventFloatName(m.dst, m.name, vel)
			}
		case core.EventMIDINoteOff:
			if me.GetNote() == m.note {
				core.SendEventFloatName(m.dst, m.name, 0)
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
