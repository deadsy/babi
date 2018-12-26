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
	info core.ModuleInfo // module info
	ch   uint8           // MIDI channel
	poly core.Module     // polyphony
	pan  core.Module     // pan left/right
}

// NewPoly returns a polyPatch module.
func NewPoly(s *core.Synth, ch uint8, sm func(s *core.Synth) core.Module) core.Module {
	log.Info.Printf("")

	const midiCtrl = 7

	// polyphony
	poly := midi.NewPoly(s, ch, sm, 16)
	// pan the output to left/right channels
	pan := mix.NewPan(s, ch, midiCtrl)

	m := &polyPatch{
		info: polyPatchInfo,
		ch:   ch,
		poly: poly,
		pan:  pan,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *polyPatch) Child() []core.Module {
	return []core.Module{m.poly, m.pan}
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
		core.EventIn(m.poly, "midi", e)
		core.EventIn(m.pan, "midi", e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *polyPatch) Process(buf ...*core.Buf) bool {
	out0 := buf[0]
	out1 := buf[1]
	// polyphony
	var out core.Buf
	m.poly.Process(&out)
	// pan left/right
	m.pan.Process(&out, out0, out1)
	return true
}

//-----------------------------------------------------------------------------
