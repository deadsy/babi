//-----------------------------------------------------------------------------
/*

Simple Patch: an ADSR envelope on a sine wave.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/audio"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/utils/log"
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

type simplePatch struct {
	synth   *core.Synth // top-level synth
	ch      uint8       // MIDI channel
	adsr    core.Module // adsr envelope
	sine    core.Module // sine oscillator
	pan     core.Module // pan left/right
	note    core.Module // note to gate
	panCtrl core.Module // MIDI to pan control
	volCtrl core.Module // MIDI to volume control
}

// NewSimplePatch returns a simple sine/adsr patch.
func NewSimplePatch(s *core.Synth) core.Module {
	log.Info.Printf("")

	const midiCh = 0
	const midiNote = 69
	const midiCtrl = 6

	adsr := env.NewADSR(s)
	sine := osc.NewSine(s)
	pan := audio.NewPan(s)
	note := midi.NewNote(s, midiCh, midiNote, adsr, "gate")
	panCtrl := midi.NewCtrl(s, midiCh, midiCtrl+0, pan, "pan")
	volCtrl := midi.NewCtrl(s, midiCh, midiCtrl+1, pan, "volume")

	// adsr defaults
	core.SendEventFloatName(adsr, "attack", 0.1)
	core.SendEventFloatName(adsr, "decay", 0.5)
	core.SendEventFloatName(adsr, "sustain", 0.1)
	core.SendEventFloatName(adsr, "release", 1)
	// sine defaults
	core.SendEventFloatName(sine, "frequency", 440.0)
	// pan defaults
	core.SendEventFloatName(pan, "pan", 0.5)
	core.SendEventFloatName(pan, "volume", 1)

	return &simplePatch{
		synth:   s,
		ch:      midiCh,
		adsr:    adsr,
		sine:    sine,
		pan:     pan,
		note:    note,
		panCtrl: panCtrl,
		volCtrl: volCtrl,
	}
}

// Return the child modules.
func (m *simplePatch) Child() []core.Module {
	return []core.Module{m.adsr, m.sine, m.pan, m.note, m.panCtrl, m.volCtrl}
}

// Stop and performs any cleanup of a module.
func (m *simplePatch) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *simplePatch) Event(e *core.Event) {
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		m.note.Event(e)
		m.panCtrl.Event(e)
		m.volCtrl.Event(e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *simplePatch) Process(buf ...*core.Buf) {
	outL := buf[0]
	outR := buf[1]
	// generate sine
	var out core.Buf
	m.sine.Process(&out)
	// generate envelope
	var env core.Buf
	m.adsr.Process(&env)
	// apply envelope
	out.Mul(&env)
	// pan left/right
	m.pan.Process(&out, outL, outR)
}

// Active return true if the module has non-zero output.
func (m *simplePatch) Active() bool {
	return m.adsr.Active()
}

//-----------------------------------------------------------------------------
