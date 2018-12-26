//-----------------------------------------------------------------------------
/*

MIDI Control Module

Convert a MIDI control message into a float event for another module.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var ctrlMidiInfo = core.ModuleInfo{
	Name: "ctrlMidi",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, ctrlMidiIn},
	},
	Out: []core.PortInfo{
		{"val", "float value (0..1)", core.PortTypeFloat, nil},
	},
}

// Info returns the module information.
func (m *ctrlMidi) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type ctrlMidi struct {
	info core.ModuleInfo // module info
	ch   uint8           // MIDI channel
	cc   uint8           // MIDI control change number
}

// NewCtrl returns a MIDI control module.
func NewCtrl(s *core.Synth, ch, cc uint8) core.Module {
	log.Info.Printf("midi ch %d cc %d", ch, cc)
	m := &ctrlMidi{
		info: ctrlMidiInfo,
		ch:   ch,
		cc:   cc,
	}
	return s.Register(m)
}

// Return the child modules.
func (m *ctrlMidi) Child() []core.Module {
	return nil
}

// Stop and cleanup the module.
func (m *ctrlMidi) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func ctrlMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*ctrlMidi)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		if me.GetType() == core.EventMIDIControlChange && me.GetCcNum() == m.cc {
			// convert to a float value and output
			core.EventOutFloat(m, "val", me.GetCcFloat())
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ctrlMidi) Process(buf ...*core.Buf) bool {
	// do nothing
	return false
}

//-----------------------------------------------------------------------------
