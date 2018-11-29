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
	return &m.info
}

//-----------------------------------------------------------------------------

type countMidi struct {
	info  core.ModuleInfo // module info
	ch    uint8           // MIDI channel
	note  uint8           // MIDI note number
	k     uint            // modulo number
	count uint            // counter
}

// NewCounter returns a MIDI counter module.
func NewCounter(s *core.Synth, ch, note uint8, k uint) core.Module {
	log.Info.Printf("")
	m := &countMidi{
		info: countMidiInfo,
		ch:   ch,
		note: note,
		k:    k,
	}
	return s.Register(m)
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
				core.EventOutInt(m, "count", int(m.count))
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
