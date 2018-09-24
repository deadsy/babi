//-----------------------------------------------------------------------------
/*

Basic Patch: an ADSR envelope on a oscillator output.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/mix"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *basicPatch) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "basic_patch",
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

type basicPatch struct {
	synth   *core.Synth // top-level synth
	ch      uint8       // MIDI channel
	osc     core.Module // oscillator
	adsr    core.Module // adsr envelope
	pan     core.Module // pan left/right
	note    core.Module // note to gate
	panCtrl core.Module // MIDI to pan control
	volCtrl core.Module // MIDI to volume control
}

// NewBasicPatch returns a basic oscillator/envelope patch.
func NewBasicPatch(s *core.Synth, osc core.Module) core.Module {
	log.Info.Printf("")

	const midiCh = 0
	const midiNote = 69
	const midiCtrl = 6

	adsr := env.NewADSREnv(s)
	pan := mix.NewPan(s)
	note := midi.NewNote(s, midiCh, midiNote, adsr, "gate")
	panCtrl := midi.NewCtrl(s, midiCh, midiCtrl+0, pan, "pan")
	volCtrl := midi.NewCtrl(s, midiCh, midiCtrl+1, pan, "volume")

	// oscillator defaults
	core.SendEventFloatName(osc, "frequency", 440.0)
	core.SendEventFloatName(osc, "duty", 0.1)
	core.SendEventFloatName(osc, "attenuation", 1.0)
	core.SendEventFloatName(osc, "slope", 0.5)
	// adsr defaults
	core.SendEventFloatName(adsr, "attack", 0.1)
	core.SendEventFloatName(adsr, "decay", 0.5)
	core.SendEventFloatName(adsr, "sustain", 0.05)
	core.SendEventFloatName(adsr, "release", 1)
	// pan defaults
	core.SendEventFloatName(pan, "pan", 0.5)
	core.SendEventFloatName(pan, "volume", 1)

	return &basicPatch{
		synth:   s,
		ch:      midiCh,
		osc:     osc,
		adsr:    adsr,
		pan:     pan,
		note:    note,
		panCtrl: panCtrl,
		volCtrl: volCtrl,
	}
}

// Return the child modules.
func (m *basicPatch) Child() []core.Module {
	return []core.Module{m.adsr, m.osc, m.pan, m.note, m.panCtrl, m.volCtrl}
}

// Stop and performs any cleanup of a module.
func (m *basicPatch) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *basicPatch) Event(e *core.Event) {
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		m.note.Event(e)
		m.panCtrl.Event(e)
		m.volCtrl.Event(e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *basicPatch) Process(buf ...*core.Buf) {
	outL := buf[0]
	outR := buf[1]
	// generate sine
	var out core.Buf
	m.osc.Process(&out)
	// generate envelope
	var env core.Buf
	m.adsr.Process(&env)
	// apply envelope
	out.Mul(&env)
	// pan left/right
	m.pan.Process(&out, outL, outR)
}

// Active return true if the module has non-zero output.
func (m *basicPatch) Active() bool {
	return m.adsr.Active()
}

//-----------------------------------------------------------------------------
