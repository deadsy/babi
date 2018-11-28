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
	Out: nil,
}

// Info returns the module information.
func (m *ctrlMidi) Info() *core.ModuleInfo {
	return &ctrlMidiInfo
}

// ID returns the unique module identifier.
func (m *ctrlMidi) ID() string {
	return m.id
}

//-----------------------------------------------------------------------------

type ctrlMidi struct {
	synth *core.Synth // top-level synth
	id    string      // module identifier
	ch    uint8       // MIDI channel
	cc    uint8       // MIDI control change number
	dst   core.Module // destination module
	name  string      // port name on destination module
}

// NewCtrl returns a MIDI control module.
func NewCtrl(s *core.Synth, ch, cc uint8, dst core.Module, name string) core.Module {
	mi := dst.Info()
	log.Info.Printf("midi ch %d cc %d controlling %s.%s port", ch, cc, mi.Name, name)
	return &ctrlMidi{
		synth: s,
		id:    core.GenerateID(ctrlMidiInfo.Name),
		ch:    ch,
		cc:    cc,
		dst:   dst,
		name:  name,
	}
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
		if me.GetType() == core.EventMIDIControlChange && me.GetCtrlNum() == m.cc {
			// convert to a float event and send
			val := core.MIDIMap(me.GetCtrlVal(), 0, 1)
			log.Info.Printf("send float event to %s.%s val %f", core.ModuleName(m.dst), m.name, val)
			core.SendEventFloat(m.dst, m.name, val)
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ctrlMidi) Process(buf ...*core.Buf) {
	// do nothing
}

// Active return true if the module has non-zero output.
func (m *ctrlMidi) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
