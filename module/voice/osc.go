//-----------------------------------------------------------------------------
/*

Oscillator Voice

This voice is a generic oscillator with an ADSR envelope applied to it.

*/
//-----------------------------------------------------------------------------

package voice

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var oscVoiceInfo = core.ModuleInfo{
	Name: "oscVoice",
	In: []core.PortInfo{
		{"gate", "oscillator gate, attack(>0) or mute(=0)", core.PortTypeFloat, oscVoiceGate},
		{"note", "midi note value", core.PortTypeFloat, oscVoiceNote},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *oscVoice) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type oscVoice struct {
	info core.ModuleInfo // module info
	adsr core.Module     // adsr envelope
	osc  core.Module     // oscillator
}

// NewOsc returns an oscillator voice module.
func NewOsc(s *core.Synth, osc core.Module) core.Module {
	log.Info.Printf("")

	adsr := env.NewADSR(s)

	// oscillator defaults
	core.EventInFloat(osc, "duty", 0.1)
	core.EventInFloat(osc, "attenuation", 1.0)
	core.EventInFloat(osc, "slope", 0.5)
	// adsr defaults
	core.EventInFloat(adsr, "attack", 0.1)
	core.EventInFloat(adsr, "decay", 0.5)
	core.EventInFloat(adsr, "sustain", 0.05)
	core.EventInFloat(adsr, "release", 1)

	m := &oscVoice{
		info: oscVoiceInfo,
		adsr: adsr,
		osc:  osc,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *oscVoice) Child() []core.Module {
	return []core.Module{m.adsr, m.osc}
}

// Stop performs any cleanup of a module.
func (m *oscVoice) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func oscVoiceGate(cm core.Module, e *core.Event) {
	m := cm.(*oscVoice)
	core.EventIn(m.adsr, "gate", e)
}

func oscVoiceNote(cm core.Module, e *core.Event) {
	m := cm.(*oscVoice)
	f := core.MIDIToFrequency(e.GetEventFloat().Val)
	core.EventInFloat(m.osc, "frequency", f)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *oscVoice) Process(buf ...*core.Buf) bool {
	// generate envelope
	var env core.Buf
	active := m.adsr.Process(&env)
	if !active {
		return false
	}
	out := buf[0]
	// generate wave
	m.osc.Process(out)
	// apply envelope
	out.Mul(&env)
	return true
}

//-----------------------------------------------------------------------------
