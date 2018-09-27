//-----------------------------------------------------------------------------
/*

Karplus Strong Patch

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/mix"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *ksPatch) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "ks",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortTypeMIDI, ksPortMidiIn},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortTypeAudioBuffer, nil},
			{"out_right", "right channel output", core.PortTypeAudioBuffer, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type ksPatch struct {
	synth   *core.Synth // top-level synth
	ch      uint8       // MIDI channel
	ks      core.Module // ks oscillator
	pan     core.Module // pan left/right
	note    core.Module // note to gate
	panCtrl core.Module // MIDI to pan control
	volCtrl core.Module // MIDI to volume control
}

// NewKarplusStrongPatch returns a karplus strong patch.
func NewKarplusStrongPatch(s *core.Synth) core.Module {
	log.Info.Printf("")

	const midiCh = 0
	const midiNote = 69
	const midiCtrl = 6

	ks := osc.NewKarplusStrong(s)
	pan := mix.NewPan(s)
	note := midi.NewNote(s, midiCh, midiNote, ks, "gate")
	panCtrl := midi.NewCtrl(s, midiCh, midiCtrl+0, pan, "pan")
	volCtrl := midi.NewCtrl(s, midiCh, midiCtrl+1, pan, "volume")

	// ks default
	core.SendEventFloat(ks, "attenuation", 1.0)
	core.SendEventFloat(ks, "frequency", 440.0)
	// pan defaults
	core.SendEventFloat(pan, "pan", 0.5)
	core.SendEventFloat(pan, "volume", 1)

	return &ksPatch{
		synth:   s,
		ch:      midiCh,
		ks:      ks,
		pan:     pan,
		note:    note,
		panCtrl: panCtrl,
		volCtrl: volCtrl,
	}
}

// Return the child modules.
func (m *ksPatch) Child() []core.Module {
	return []core.Module{m.ks, m.pan, m.note, m.panCtrl, m.volCtrl}
}

// Stop and performs any cleanup of a module.
func (m *ksPatch) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Port Events

func ksPortMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*ksPatch)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		core.SendEvent(m.note, "midi_in", e)
		core.SendEvent(m.panCtrl, "midi_in", e)
		core.SendEvent(m.volCtrl, "midi_in", e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ksPatch) Process(buf ...*core.Buf) {
	outL := buf[0]
	outR := buf[1]
	// generate wave
	var out core.Buf
	m.ks.Process(&out)
	// pan left/right
	m.pan.Process(&out, outL, outR)
}

// Active return true if the module has non-zero output.
func (m *ksPatch) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
