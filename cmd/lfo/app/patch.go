//-----------------------------------------------------------------------------
/*

LFO test patch.

*/
//-----------------------------------------------------------------------------

package app

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/module/mix"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var patchAppInfo = core.ModuleInfo{
	Name: "patchApp",
	In: []core.PortInfo{
		{"midi", "midi input", core.PortTypeMIDI, patchAppMidiIn},
	},
	Out: []core.PortInfo{
		{"out0", "left channel output", core.PortTypeAudio, nil},
		{"out1", "right channel output", core.PortTypeAudio, nil},
	},
}

// Info returns the general module information.
func (m *patchApp) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type patchApp struct {
	info core.ModuleInfo // module info
	ctrl core.Module     // MIDI filter/processor
	poly core.Module     // polyphony
	pan  core.Module     // pan left/right
}

// NewPatch returns an LFO test patch.
func NewPatch(s *core.Synth, ch uint8) core.Module {

	// process incoming midi
	ctrl := NewCtrl(s, ch)

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
	m := &patchApp{
		info: patchAppInfo,
		ctrl: ctrl,
		poly: poly,
		pan:  pan,
	}

	// set the initial cc values
	core.EventInBool(ctrl, "reset", true)

	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *patchApp) Child() []core.Module {
	return []core.Module{m.ctrl, m.poly, m.pan}
}

// Stop performs any cleanup of a module.
func (m *patchApp) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func patchAppMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*patchApp)
	core.EventIn(m.ctrl, "midi", e)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *patchApp) Process(buf ...*core.Buf) {
	out0 := buf[0]
	out1 := buf[1]
	// polyphony
	var out core.Buf
	m.poly.Process(&out)
	// pan left/right
	m.pan.Process(&out, out0, out1)
}

// Active returns true if the module has non-zero output.
func (m *patchApp) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
