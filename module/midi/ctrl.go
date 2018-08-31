//-----------------------------------------------------------------------------
/*

MIDI Control Module

Convert a MIDI control message into a float event for another module.

*/
//-----------------------------------------------------------------------------

package midi

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *ctrlModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "midi_control",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortType_EventMIDI, 0},
		},
		Out: nil,
	}
}

//-----------------------------------------------------------------------------

type ctrlModule struct {
	ch   uint8       // MIDI channel
	cc   uint8       // MIDI control change number
	dst  core.Module // destination module
	ctrl uint        // control number for the destination module
}

// NewCtrl returns a MIDI control module.
func NewCtrl(ch, cc uint8, dst core.Module, name string) core.Module {
	log.Info.Printf("")
	return &ctrlModule{
		ch:   ch,
		cc:   cc,
		dst:  dst,
		ctrl: dst.Info().GetPortID(name),
	}
}

// Stop and cleanup the module.
func (m *ctrlModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------

// Event processes a module event.
func (m *ctrlModule) Event(e *core.Event) {
	log.Info.Printf("event %s", e)
	if e.IsMIDI(m.ch) {
		me := e.GetEventMIDI()
		if me.GetType() == core.EventMIDI_ControlChange && me.GetCtrlNum() == m.cc {
			// convert to a float event and send
			val := core.MIDI_Map(me.GetCtrlVal(), 0, 1)
			log.Info.Printf("send float event to %s port %d val %f", m.dst, m.ctrl, val)
			m.dst.Event(core.NewEventFloat(m.ctrl, val))
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ctrlModule) Process(buf ...*core.Buf) {
	// do nothing
}

// Active return true if the module has non-zero output.
func (m *ctrlModule) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
