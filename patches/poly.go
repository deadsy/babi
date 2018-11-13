//-----------------------------------------------------------------------------
/*

Poly Patch

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/mix"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *polyPatch) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "polyPatch",
		In: []core.PortInfo{
			{"midi_in", "midi input", core.PortTypeMIDI, polyPatchMidiIn},
		},
		Out: []core.PortInfo{
			{"out_left", "left channel output", core.PortTypeAudioBuffer, nil},
			{"out_right", "right channel output", core.PortTypeAudioBuffer, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type polyPatch struct {
	synth   *core.Synth // top-level synth
	ch      uint8       // MIDI channel
	poly    core.Module // polyphony
	pan     core.Module // pan left/right
	panCtrl core.Module // MIDI to pan control
	volCtrl core.Module // MIDI to volume control
}

// NewPoly returns a polyPatch module.
func NewPoly(s *core.Synth, sm func(s *core.Synth) core.Module) core.Module {
	log.Info.Printf("")

	const midiCh = 0
	const midiCtrl = 6

	// polyphony
	poly := midi.NewPoly(s, midiCh, sm, 16)
	// pan the output to left/right channels
	pan := mix.NewPan(s)
	panCtrl := midi.NewCtrl(s, midiCh, midiCtrl+0, pan, "pan")
	volCtrl := midi.NewCtrl(s, midiCh, midiCtrl+1, pan, "volume")

	// pan defaults
	core.SendEventFloat(pan, "pan", 0.5)
	core.SendEventFloat(pan, "volume", 1)

	return &polyPatch{
		synth:   s,
		ch:      midiCh,
		poly:    poly,
		pan:     pan,
		panCtrl: panCtrl,
		volCtrl: volCtrl,
	}
}

// Child returns the child modules of this module.
func (m *polyPatch) Child() []core.Module {
	return []core.Module{m.poly, m.pan, m.panCtrl, m.volCtrl}
}

// Stop performs any cleanup of a module.
func (m *polyPatch) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func polyPatchMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*polyPatch)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		core.SendEvent(m.poly, "midi_in", e)
		core.SendEvent(m.panCtrl, "midi_in", e)
		core.SendEvent(m.volCtrl, "midi_in", e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *polyPatch) Process(buf ...*core.Buf) {
	outL := buf[0]
	outR := buf[1]
	// polyphony
	var out core.Buf
	m.poly.Process(&out)
	// pan left/right
	m.pan.Process(&out, outL, outR)
}

// Active returns true if the module has non-zero output.
func (m *polyPatch) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
