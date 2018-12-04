//-----------------------------------------------------------------------------
/*

Goom Voice Control Module

This is a MIDI event processor.

A goom voice has 22 controls and 2 modal switches.
My MIDI controller (AKAI MPKmini) has 8 CC controls.
We use a drum pad note as a modal switch to multiplex CC controls into 3 groups.

Group 0 (wave):
duty slope pan vol
attack decay sustain release

Group 1 (modulation):
duty slope X level
attack decay coarse fine

Group 2 (filter):
sensitivity cutoff X resonance
attack decay sustain release

*/
//-----------------------------------------------------------------------------

package goom

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const midiOscillatorModeNote = 45 // note for oscillator mode
const midiFrequencyModeNote = 46  // note for frequency mode
const midiCCModeNote = 47         // note for cc mode

const nControls = 8 // cc controls per mode

// 3 x 8 CC values
const midiWaveDutyCC = 1           // wave oscillator duty cycle
const midiWaveSlopeCC = 2          // wave oscillator duty slope
const midiPanCC = 3                // pan left/right
const midiPanVolCC = midiPanCC + 1 // main volume
const midiAmpAttackCC = 5          // amplitude attack
const midiAmpDecayCC = 6           // amplitude decay
const midiAmpSustainCC = 7         // amplitude sustain
const midiAmpReleaseCC = 8         // amplitude release
const midiUnusedCC9 = 9            // unused
const midiUnusedCC10 = 10          // unused
const midiUnusedCC11 = 11          // unused
const midiUnusedCC12 = 12          // unused
const midiUnusedCC13 = 13          // unused
const midiUnusedCC14 = 14          // unused
const midiUnusedCC15 = 15          // unused
const midiUnusedCC16 = 16          // unused
const midiUnusedCC17 = 17          // unused
const midiUnusedCC18 = 18          // unused
const midiUnusedCC19 = 19          // unused
const midiUnusedCC20 = 20          // unused
const midiUnusedCC21 = 21          // unused
const midiUnusedCC22 = 22          // unused
const midiUnusedC23 = 23           // unused
const midiUnusedCC24 = 24          // unused

// and keys turned into CCs
const midiOscillatorModeCC = 25 // oscillator mode cc
const midiFrequencyModeCC = 26  // frequency mode cc

//-----------------------------------------------------------------------------

var ctrlGoomInfo = core.ModuleInfo{
	Name: "ctrlGoom",
	In: []core.PortInfo{
		{"midi", "midi in", core.PortTypeMIDI, ctrlGoomMidiIn},
		{"reset", "reset cc values", core.PortTypeBool, ctrlGoomReset},
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

func ctrlGoomReset(cm core.Module, e *core.Event) {
	m := cm.(*ctrlGoom)
	be := e.GetEventBool()
	if be != nil && be.Val {
		log.Info.Printf("")
		core.EventOutMidiCC(m, "midi", midiWaveDutyCC, 64)
		core.EventOutMidiCC(m, "midi", midiWaveSlopeCC, 64)
		// amplitude envelope
		core.EventOutMidiCC(m, "midi", midiAmpAttackCC, 64)
		core.EventOutMidiCC(m, "midi", midiAmpDecayCC, 64)
		core.EventOutMidiCC(m, "midi", midiAmpSustainCC, 64)
		core.EventOutMidiCC(m, "midi", midiAmpReleaseCC, 64)
		// output mixing
		core.EventOutMidiCC(m, "midi", midiPanCC, 64)
		core.EventOutMidiCC(m, "midi", midiPanVolCC, 64)
		// oscillator/frequency modes
		core.EventOutMidiCC(m, "midi", midiOscillatorModeCC, 0)
		core.EventOutMidiCC(m, "midi", midiFrequencyModeCC, 0)
	}
}

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
			ccNum := me.GetCcNum()
			if ccNum >= 1 && ccNum <= 8 {
				ccNum += m.ccMode * nControls
				core.EventOutMidiCC(m, "midi", ccNum, me.GetCcInt())
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
