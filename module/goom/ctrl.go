//-----------------------------------------------------------------------------
/*

Goom Voice Control Module

A goom voice has more controls than I have knobs on my AKAI MPKmini MIDI controller.
This module alows modal switching between the 8 knobs I do have.
That is: Hit a drum pad, switch modes to a different control group.

*/
//-----------------------------------------------------------------------------

package goom

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/midi"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const midiChannel = 0 // midi channel

const midiNote0 = 49 // note for oscillator mode control
const midiNote1 = 50 // note for frequency mode control
const midiNote2 = 51 // note for control channel mode switching
const midiCC = 1     // base control
const nControls = 8  // controls per mode
const nModes = 3     // number of modes

//-----------------------------------------------------------------------------

var ctrlGoomInfo = core.ModuleInfo{
	Name: "ctrlGoom",
	In: []core.PortInfo{
		{"midi", "midi", core.PortTypeMIDI, ctrlGoomMidiIn},
	},
	Out: nil,
}

// Info returns the module information.
func (m *ctrlGoom) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type ctrlGoom struct {
	info    core.ModuleInfo // module info
	oMode   core.Module     // oscillator mode sleection
	fMode   core.Module     // frequency mode selection
	ccMode  core.Module     // control channel mode selection
	ccModal core.Module     //
}

// NewController returns a goom voice MIDI controller.
func NewController(s *core.Synth) core.Module {
	log.Info.Printf("")

	oMode := midi.NewCounter(s, midiChannel, midiNote0, 3)
	fMode := midi.NewCounter(s, midiChannel, midiNote1, 3)
	ccMode := midi.NewCounter(s, midiChannel, midiNote2, nModes)
	ccModal := midi.NewModal(s, midiChannel, midiCC, nControls, nModes)

	core.Connect(ccMode, "count", ccModal, "mode")

	m := &ctrlGoom{
		info:    ctrlGoomInfo,
		oMode:   oMode,
		fMode:   fMode,
		ccMode:  ccMode,
		ccModal: ccModal,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *ctrlGoom) Child() []core.Module {
	return []core.Module{m.oMode, m.fMode, m.ccMode, m.ccModal}
}

// Stop performs any cleanup of a module.
func (m *ctrlGoom) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func ctrlGoomMidiIn(cm core.Module, e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ctrlGoom) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *ctrlGoom) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
