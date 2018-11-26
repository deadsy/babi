//-----------------------------------------------------------------------------
/*

MIDI Note Counter Module

Increment a modulo counter with every note on event.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *intMidi) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "intMidi",
		In: []core.PortInfo{
			{"midi", "midi", core.PortTypeMIDI, intMidiIn},
		},
		Out: []core.PortInfo{
			{"n", "counter", core.PortTypeInt, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type intMidi struct {
	synth *core.Synth // top-level synth
	ch    uint8       // MIDI channel
	note  uint8       // MIDI note number
	k     uint        // modulo number
	count uint        // counter
}

// NewIntMidi returns a MIDI counter module.
func NewIntMidi(s *core.Synth, ch, note uint8, k uint) core.Module {
	log.Info.Printf("")
	return &intMidi{
		synth: s,
		ch:    ch,
		note:  note,
		k:     k,
	}
}

// Child returns the child modules of this module.
func (m *intMidi) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *intMidi) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func intMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*intMidi)

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
func (m *intMidi) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *intMidi) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
