//-----------------------------------------------------------------------------
/*

Simple Patch: an ADSR envelope on a sine wave.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
	"github.com/deadsy/babi/module/audio"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/osc"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *simplePatch) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "simple_patch",
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

const midi_ch = 0
const midi_note = 69
const midi_ctrl = 6

type simplePatch struct {
	ch       uint8       // MIDI channel
	adsr     core.Module // adsr envelope
	sine     core.Module // sine oscillator
	pan      core.Module // pan left/right
	note     core.Module // note to gate
	pan_ctrl core.Module // MIDI to pan control
	vol_ctrl core.Module // MIDI to volume control
}

func NewSimple() core.Module {
	log.Info.Printf("")

	adsr := env.NewADSR()
	sine := osc.NewSine()
	pan := audio.NewPan()
	note := midi.NewNote(midi_ch, midi_note, adsr, "gate")
	pan_ctrl := midi.NewCtrl(midi_ch, midi_ctrl+0, pan, "pan")
	vol_ctrl := midi.NewCtrl(midi_ch, midi_ctrl+1, pan, "volume")

	// adsr defaults
	adsr.Event(core.NewEventFloat(adsr.Info().GetPortByName("attack").Id, 0.1))
	adsr.Event(core.NewEventFloat(adsr.Info().GetPortByName("decay").Id, 0.5))
	adsr.Event(core.NewEventFloat(adsr.Info().GetPortByName("sustain").Id, 0.7))
	adsr.Event(core.NewEventFloat(adsr.Info().GetPortByName("release").Id, 1))
	// sine defaults
	sine.Event(core.NewEventFloat(sine.Info().GetPortByName("frequency").Id, 440.0))
	// pan defaults
	pan.Event(core.NewEventFloat(pan.Info().GetPortByName("pan").Id, 0.5))
	pan.Event(core.NewEventFloat(pan.Info().GetPortByName("volume").Id, 1))

	return &simplePatch{
		ch:       midi_ch,
		adsr:     adsr,
		sine:     sine,
		pan:      pan,
		note:     note,
		pan_ctrl: pan_ctrl,
		vol_ctrl: vol_ctrl,
	}
}

// Stop and performs any cleanup of a module.
func (m *simplePatch) Stop() {
	log.Info.Printf("")
	m.adsr.Stop()
	m.sine.Stop()
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *simplePatch) Event(e *core.Event) {
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		m.note.Event(e)
		m.pan_ctrl.Event(e)
		m.vol_ctrl.Event(e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *simplePatch) Process(buf ...*core.Buf) {
	out_l := buf[0]
	out_r := buf[1]
	// generate sine
	var out core.Buf
	m.sine.Process(&out)
	// generate envelope
	var env core.Buf
	m.adsr.Process(&env)
	// apply envelope
	out.Mul(&env)
	// pan left/right
	m.pan.Process(&out, out_l, out_r)
}

// Active return true if the module has non-zero output.
func (m *simplePatch) Active() bool {
	return m.adsr.Active()
}

//-----------------------------------------------------------------------------
