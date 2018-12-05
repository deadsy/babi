//-----------------------------------------------------------------------------
/*

LFO Test Control Module

This is a MIDI event processor.

*/
//-----------------------------------------------------------------------------

package app

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const midiLfoModeNote = 47 // note for lfo mode

const midiWaveDutyCC = 1           // wave oscillator duty cycle
const midiWaveSlopeCC = 2          // wave oscillator duty slope
const midiPanCC = 3                // pan left/right
const midiPanVolCC = midiPanCC + 1 // main volume
const midiAmpAttackCC = 5          // amplitude attack
const midiAmpDecayCC = 6           // amplitude decay
const midiAmpSustainCC = 7         // amplitude sustain
const midiAmpReleaseCC = 8         // amplitude release

// and keys turned into CCs
const midiLfoModeCC = 25 // lfo mode

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
	info    core.ModuleInfo // module info
	ch      uint8           // MIDI channel
	lfoMode uint8           // lfo mode (0..5)
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
		// lfo mode
		core.EventOutMidiCC(m, "midi", midiLfoModeCC, 0)
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
			case midiLfoModeNote:
				m.lfoMode = (m.lfoMode + 1) % 6
				log.Info.Printf("lfoMode %d", m.lfoMode)
				core.EventOutMidiCC(m, "midi", midiLfoModeCC, m.lfoMode)
				return
			}
		// Ignore the note off events for our special keys.
		case core.EventMIDINoteOff:
			switch me.GetNote() {
			case midiLfoModeNote:
				// filter out
				return
			}
		}
		// pass through
		core.EventOut(m, "midi", e)
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ctrlApp) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *ctrlApp) Active() bool {
	return false
}

//-----------------------------------------------------------------------------
