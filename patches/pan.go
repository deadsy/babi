//-----------------------------------------------------------------------------
/*

Pan Patch:

Single output channel module with a pan module added for L/R output.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
	"github.com/deadsy/babi/module/audio"
	"github.com/deadsy/babi/module/midi"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *panPatch) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "pan_patch",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortType_EventMIDI, 0},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortType_AudioBuffer, 0},
			{"out_right", "right channel output", core.PortType_AudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type panPatch struct {
	ch       uint8       // MIDI channel
	sm       core.Module // sub-module
	pan      core.Module // pan module
	pan_ctrl core.Module // MIDI control for pan
	vol_ctrl core.Module // MIDI control for volume
}

func NewPan(ch uint8, sm core.Module) core.Module {
	log.Info.Printf("")
	pan := audio.NewPan()
	return &panPatch{
		ch:       ch,
		sm:       sm,
		pan:      pan,
		pan_ctrl: midi.NewCtrl(ch, 10, pan, "pan"),
		vol_ctrl: midi.NewCtrl(ch, 11, pan, "volume"),
	}
}

// Stop and cleanup the module.
func (m *panPatch) Stop() {
	log.Info.Printf("")
	m.sm.Stop()
	m.pan.Stop()
}

//-----------------------------------------------------------------------------

// Event processes a module event.
func (m *panPatch) Event(e *core.Event) {
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		m.sm.Event(e)
		m.pan_ctrl.Event(e)
		m.vol_ctrl.Event(e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *panPatch) Process(buf ...*core.Buf) {
	out_l := buf[0]
	out_r := buf[0]
	var out core.Buf
	m.sm.Process(&out)
	m.pan.Process(&out, out_l, out_r)
}

// Active return true if the module has non-zero output.
func (m *panPatch) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
