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
			{"midi", "midi input", core.PortType_EventMIDI, 0},
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

// Stop and performs any cleanup of a module.
func (m *ctrlModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------

// Event processes a module event.
func (m *ctrlModule) Event(e *core.Event) {
	log.Info.Printf("event %s", e)
	switch e.GetType() {
	case core.Event_MIDI:
		me := e.GetEventMIDI()
		switch me.GetType() {
		case core.EventMIDI_ControlChange:
			// filter on channel and control number
			if me.GetChannel() == m.ch && me.GetCtrlNum() == m.cc {
				// convert to a float event and send
				val := core.MIDI_Map(me.GetCtrlVal(), 0, 1)
				m.dst.Event(core.NewEventFloat(m.ctrl, val))
				return
			}
			fallthrough
		default:
			log.Info.Printf("unhandled midi event %s", me)
		}
	default:
		log.Info.Printf("unhandled event %s", e)
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
