//-----------------------------------------------------------------------------
/*

Goom synth root level patch.

*/
//-----------------------------------------------------------------------------

package goom

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/mix"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var patchGoomInfo = core.ModuleInfo{
	Name: "patchGoom",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, patchGoomMidiIn},
	},
	Out: []core.PortInfo{
		{"out0", "left channel output", core.PortTypeAudio, nil},
		{"out1", "right channel output", core.PortTypeAudio, nil},
	},
}

// Info returns the general module information.
func (m *patchGoom) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type patchGoom struct {
	info core.ModuleInfo // module info
	ctrl core.Module     // MIDI filter/processor
	poly core.Module     // polyphony
	pan  core.Module     // pan left/right
}

// NewPatch returns an goom root module.
func NewPatch(s *core.Synth, ch uint8) core.Module {

	// process incoming midi
	ctrl := NewCtrl(s, ch)

	// function to create a goom voices
	//voice := func(s *core.Synth) core.Module { return NewVoice(s) }

	// polyphony
	poly := midi.NewPoly(s, ch, NewVoice, 16)
	core.Connect(ctrl, "midi", poly, "midi")

	// pan the output to left/right channels
	pan := mix.NewPan(s, ch, midiPanCC)
	core.Connect(ctrl, "midi", pan, "midi")

	// monitor the MIDI events
	mon := midi.NewMonitor(s, ch)
	core.Connect(ctrl, "midi", mon, "midi")

	log.Info.Printf("")
	m := &patchGoom{
		info: patchGoomInfo,
		ctrl: ctrl,
		poly: poly,
		pan:  pan,
	}

	// set the initial cc values
	core.EventInBool(ctrl, "reset", true)

	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *patchGoom) Child() []core.Module {
	return []core.Module{m.ctrl, m.poly, m.pan}
}

// Stop performs any cleanup of a module.
func (m *patchGoom) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func patchGoomMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*patchGoom)
	core.EventIn(m.ctrl, "midi", e)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *patchGoom) Process(buf ...*core.Buf) bool {
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
