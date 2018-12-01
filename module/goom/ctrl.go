//-----------------------------------------------------------------------------
/*

Goom Voice Control Module

A goom voice has about 21 controls.
My MIDI controller (AKAI MPKmini) has 8 CC knobs.
This MIDI event processor uses drum pads as modal controls to multiplex
the CC controls into multiple groups.

*/
//-----------------------------------------------------------------------------

package goom

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const midiOscillatorModeNote = 49 // note for oscillator mode
const midiFrequencyModeNote = 50  // note for frequency mode
const midiCCModeNote = 51         // note for cc mode

const midiPanCC = 23            // pan and vol(+1) cc
const midiOscillatorModeCC = 25 // oscillator mode cc
const midiFrequencyModeCC = 26  // frequency mode cc

const nControls = 8 // cc controls per mode

//-----------------------------------------------------------------------------

var ctrlGoomInfo = core.ModuleInfo{
	Name: "ctrlGoom",
	In: []core.PortInfo{
		{"midi", "midi in", core.PortTypeMIDI, ctrlGoomMidiIn},
	},
	Out: []core.PortInfo{
		{"midi", "midi out", core.PortTypeMIDI, nil},
	},
}

// Info returns the module information.
func (m *ctrlGoom) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type ctrlGoom struct {
	info   core.ModuleInfo // module info
	ch     uint8           // MIDI channel
	oMode  uint8           // oscillator mode (0,1,2)
	fMode  uint8           // frequency mode (0,1,2)
	ccMode uint8           // cc mode  (0,1,2)
}

// NewCtrl returns a goom voice MIDI controller.
func NewCtrl(s *core.Synth, ch uint8) core.Module {
	log.Info.Printf("")
	m := &ctrlGoom{
		info: ctrlGoomInfo,
		ch:   ch,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *ctrlGoom) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *ctrlGoom) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func ctrlGoomMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*ctrlGoom)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		// Use the special key note on events to modulo increment the mode variables.
		case core.EventMIDINoteOn:
			switch me.GetNote() {
			case midiOscillatorModeNote:
				m.oMode = (m.oMode + 1) % 3
				log.Info.Printf("omode %d", m.oMode)
				core.EventOutMidiCC(m, "midi", midiOscillatorModeCC, m.oMode)
				return
			case midiFrequencyModeNote:
				m.fMode = (m.fMode + 1) % 3
				log.Info.Printf("fmode %d", m.fMode)
				core.EventOutMidiCC(m, "midi", midiFrequencyModeCC, m.fMode)
				return
			case midiCCModeNote:
				m.ccMode = (m.ccMode + 1) % 3
				log.Info.Printf("ccmode %d", m.ccMode)
				return
			}
		// Ignore the note off events for our special keys.
		case core.EventMIDINoteOff:
			switch me.GetNote() {
			case midiOscillatorModeNote,
				midiFrequencyModeNote,
				midiCCModeNote:
				// filter out
				return
			}
		// Re-emit the CC events with higher CC numbers.
		case core.EventMIDIControlChange:
			ccNum := me.GetCtrlNum()
			if ccNum >= 1 && ccNum <= 8 {
				ccNum += m.ccMode * nControls
				core.EventOutMidiCC(m, "midi", ccNum, me.GetCtrlVal())
				return
			}
		}
		// pass through
		core.EventOut(m, "midi", e)
	}
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
