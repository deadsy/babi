//-----------------------------------------------------------------------------
/*

Poly Patch

*/
//-----------------------------------------------------------------------------

package patch

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/mix"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var polyPatchInfo = core.ModuleInfo{
	Name: "polyPatch",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, polyPatchMidiIn},
	},
	Out: []core.PortInfo{
		{"out0", "left channel output", core.PortTypeAudio, nil},
		{"out1", "right channel output", core.PortTypeAudio, nil},
	},
}

// Info returns the general module information.
func (m *polyPatch) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type polyPatch struct {
	info    core.ModuleInfo // module info
	ch      uint8           // MIDI channel
	poly    core.Module     // polyphony
	pan     core.Module     // pan left/right
	panCtrl core.Module     // MIDI to pan control
	volCtrl core.Module     // MIDI to volume control
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
	panCtrl := midi.NewCtrl(s, midiCh, midiCtrl+0)
	core.Connect(panCtrl, "val", pan, "pan")
	volCtrl := midi.NewCtrl(s, midiCh, midiCtrl+1)
	core.Connect(volCtrl, "val", pan, "vol")

	// pan defaults
	core.SendEventFloat(pan, "pan", 0.5)
	core.SendEventFloat(pan, "vol", 1)

	m := &polyPatch{
		info:    polyPatchInfo,
		ch:      midiCh,
		poly:    poly,
		pan:     pan,
		panCtrl: panCtrl,
		volCtrl: volCtrl,
	}
	return s.Register(m)
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
		core.SendEvent(m.poly, "midi", e)
		core.SendEvent(m.panCtrl, "midi", e)
		core.SendEvent(m.volCtrl, "midi", e)
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
