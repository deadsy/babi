//-----------------------------------------------------------------------------
/*

Pan Patch:

Single output channel module with a pan module added for L/R output.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/audio"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/utils/log"
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
	ch      uint8       // MIDI channel
	sm      core.Module // sub-module
	pan     core.Module // pan module
	panCtrl core.Module // MIDI control for pan
	volCtrl core.Module // MIDI control for volume
}

// NewPan returns a module that pans the output of a given sub-module.
func NewPan(ch uint8, sm core.Module) core.Module {
	log.Info.Printf("")
	// check for IO compatibility
	err := sm.Info().CheckIO(1, 0, 1)
	if err != nil {
		panic(err)
	}
	pan := audio.NewPan()
	return &panPatch{
		ch:      ch,
		sm:      sm,
		pan:     pan,
		panCtrl: midi.NewCtrl(ch, 10, pan, "pan"),
		volCtrl: midi.NewCtrl(ch, 11, pan, "volume"),
	}
}

// Return the child modules.
func (m *panPatch) Child() []core.Module {
	return []core.Module{m.sm, m.pan, m.panCtrl, m.volCtrl}
}

// Stop and cleanup the module.
func (m *panPatch) Stop() {
	log.Info.Printf("")
	m.pan.Stop()
	m.panCtrl.Stop()
	m.volCtrl.Stop()
}

//-----------------------------------------------------------------------------

// Event processes a module event.
func (m *panPatch) Event(e *core.Event) {
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		m.sm.Event(e)
		m.panCtrl.Event(e)
		m.volCtrl.Event(e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *panPatch) Process(buf ...*core.Buf) {
	outL := buf[0]
	outR := buf[0]
	var out core.Buf
	m.sm.Process(&out)
	m.pan.Process(&out, outL, outR)
}

// Active return true if the module has non-zero output.
func (m *panPatch) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
