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

var noteMidiInfo = core.ModuleInfo{
	Name: "noteMidi",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, noteMidiIn},
	},
	Out: []core.PortInfo{
		{"gate", "gate ouput", core.PortTypeFloat, nil},
	},
}

// Info returns the module information.
func (m *noteMidi) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type noteMidi struct {
	info core.ModuleInfo // module info
	ch   uint8           // MIDI channel
	note uint8           // MIDI note number
}

// NewNote returns a MIDI note on/off gate control module.
func NewNote(s *core.Synth, ch, note uint8) core.Module {
	log.Info.Printf("midi ch %d note %d", ch, note)
	m := &noteMidi{
		info: noteMidiInfo,
		ch:   ch,
		note: note,
	}
	return s.Register(m)
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
				core.EventOutFloat(m, "gate", vel)
			}
		case core.EventMIDINoteOff:
			if me.GetNote() == m.note {
				core.EventOutFloat(m, "gate", 0)
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

// Active returns true if the module has non-zero output.
func (m *noteMidi) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
