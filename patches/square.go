//-----------------------------------------------------------------------------
/*

Square Patch: an ADSR envelope on a square wave.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/mix"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *sqrPatch) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "square_patch",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortTypeMIDI, 0},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortTypeAudioBuffer, 0},
			{"out_right", "right channel output", core.PortTypeAudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type sqrPatch struct {
	synth   *core.Synth // top-level synth
	ch      uint8       // MIDI channel
	adsr    core.Module // adsr envelope
	sqr     core.Module // square oscillator
	pan     core.Module // pan left/right
	note    core.Module // note to gate
	panCtrl core.Module // MIDI to pan control
	volCtrl core.Module // MIDI to volume control
}

// NewSquarePatch returns a simple square/adsr patch.
func NewSquarePatch(s *core.Synth) core.Module {
	log.Info.Printf("")

	const midiCh = 0
	const midiNote = 69
	const midiCtrl = 6

	adsr := env.NewADSR(s)
	sqr := osc.NewSquareBasic(s)
	pan := mix.NewPan(s)
	note := midi.NewNote(s, midiCh, midiNote, adsr, "gate")
	panCtrl := midi.NewCtrl(s, midiCh, midiCtrl+0, pan, "pan")
	volCtrl := midi.NewCtrl(s, midiCh, midiCtrl+1, pan, "volume")

	// adsr defaults
	core.SendEventFloatName(adsr, "attack", 0.1)
	core.SendEventFloatName(adsr, "decay", 0.5)
	core.SendEventFloatName(adsr, "sustain", 0.05)
	core.SendEventFloatName(adsr, "release", 1)
	// sine defaults
	core.SendEventFloatName(sqr, "frequency", 440.0)
	core.SendEventFloatName(sqr, "duty", 0.1)
	// pan defaults
	core.SendEventFloatName(pan, "pan", 0.5)
	core.SendEventFloatName(pan, "volume", 1)

	return &sqrPatch{
		synth:   s,
		ch:      midiCh,
		adsr:    adsr,
		sqr:     sqr,
		pan:     pan,
		note:    note,
		panCtrl: panCtrl,
		volCtrl: volCtrl,
	}
}

// Return the child modules.
func (m *sqrPatch) Child() []core.Module {
	return []core.Module{m.adsr, m.sqr, m.pan, m.note, m.panCtrl, m.volCtrl}
}

// Stop and performs any cleanup of a module.
func (m *sqrPatch) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *sqrPatch) Event(e *core.Event) {
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		m.note.Event(e)
		m.panCtrl.Event(e)
		m.volCtrl.Event(e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *sqrPatch) Process(buf ...*core.Buf) {
	outL := buf[0]
	outR := buf[1]
	// generate sine
	var out core.Buf
	m.sqr.Process(&out)
	// generate envelope
	var env core.Buf
	m.adsr.Process(&env)
	// apply envelope
	out.Mul(&env)
	// pan left/right
	m.pan.Process(&out, outL, outR)
}

// Active return true if the module has non-zero output.
func (m *sqrPatch) Active() bool {
	return m.adsr.Active()
}

//-----------------------------------------------------------------------------
