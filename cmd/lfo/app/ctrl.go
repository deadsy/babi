//-----------------------------------------------------------------------------
/*

LFO Test Control Module

This is a MIDI event processor.

*/
//-----------------------------------------------------------------------------

package app

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const nControls = 8 // cc controls per mode

const midiModModeNote = 45  // note for modulation mode (am/fm/pm)
const midiLfoShapeNote = 46 // note for LFO wave shape
const midiCCModeNote = 47   // note for CC mode

// 1..8
const midiWaveDutyCC = 1           // wave oscillator duty cycle
const midiWaveSlopeCC = 2          // wave oscillator duty slope
const midiPanCC = 3                // pan left/right
const midiPanVolCC = midiPanCC + 1 // main volume
const midiAmpAttackCC = 5          // amplitude attack
const midiAmpDecayCC = 6           // amplitude decay
const midiAmpSustainCC = 7         // amplitude sustain
const midiAmpReleaseCC = 8         // amplitude release

// 9..16
const midiLfoRateCC = 9   // lfo rate
const midiLfoDepthCC = 10 // lfo depth

// and keys turned into CCs
const midiLfoShapeCC = 25 // lfo shape
const midiModModeCC = 26  // modulation mode (am/fm/pm)

//-----------------------------------------------------------------------------

var ctrlAppInfo = core.ModuleInfo{
	Name: "ctrlApp",
	In: []core.PortInfo{
		{"midi", "midi in", core.PortTypeMIDI, ctrlAppMidiIn},
		{"reset", "reset cc values", core.PortTypeBool, ctrlAppReset},
	},
	Out: []core.PortInfo{
		{"midi", "midi out", core.PortTypeMIDI, nil},
	},
}

// Info returns the module information.
func (m *ctrlApp) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type ctrlApp struct {
	info     core.ModuleInfo // module info
	ch       uint8           // MIDI channel
	lfoShape uint8           // lfo shape
	modMode  uint8           // modulation mode
	ccMode   uint8           // cc mode  (0,1)
}

// NewCtrl returns a goom voice MIDI controller.
func NewCtrl(s *core.Synth, ch uint8) core.Module {
	log.Info.Printf("")
	m := &ctrlApp{
		info: ctrlAppInfo,
		ch:   ch,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *ctrlApp) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *ctrlApp) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func ctrlAppReset(cm core.Module, e *core.Event) {
	m := cm.(*ctrlApp)
	be := e.GetEventBool()
	if be != nil && be.Val {
		log.Info.Printf("")
		// wave shape
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
		// lfo
		core.EventOutMidiCC(m, "midi", midiLfoRateCC, 64)
		core.EventOutMidiCC(m, "midi", midiLfoDepthCC, 64)
		core.EventOutMidiCC(m, "midi", midiModModeCC, 2)  // fm
		core.EventOutMidiCC(m, "midi", midiLfoShapeCC, 0) // triangle
	}
}

func ctrlAppMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*ctrlApp)
	me := e.GetEventMIDIChannel(m.ch)
	if me != nil {
		switch me.GetType() {
		// Use the special key note on events to modulo increment the mode variables.
		case core.EventMIDINoteOn:
			switch me.GetNote() {
			case midiModModeNote:
				m.modMode = (m.modMode + 1) % 4
				log.Info.Printf("modulation mode %s", modMode(m.modMode))
				core.EventOutMidiCC(m, "midi", midiModModeCC, m.modMode)
				return
			case midiLfoShapeNote:
				m.lfoShape = (m.lfoShape + 1) % 6
				log.Info.Printf("lfo shape %s", osc.LfoWaveShape(m.lfoShape))
				core.EventOutMidiCC(m, "midi", midiLfoShapeCC, m.lfoShape)
				return
			case midiCCModeNote:
				m.ccMode = (m.ccMode + 1) % 2
				log.Info.Printf("ccmode %d", m.ccMode)
				return
			}
		// Ignore the note off events for our special keys.
		case core.EventMIDINoteOff:
			switch me.GetNote() {
			case midiModModeNote,
				midiLfoShapeNote,
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
func (m *ctrlApp) Process(buf ...*core.Buf) bool {
	return false
}

//-----------------------------------------------------------------------------
