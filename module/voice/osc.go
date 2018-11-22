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

// Info returns the module information.
func (m *oscVoice) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "oscVoice",
		In: []core.PortInfo{
			{"gate", "oscillator gate, attack(>0) or mute(=0)", core.PortTypeFloat, oscVoiceGate},
			{"frequency", "frequency (Hz)", core.PortTypeFloat, oscVoiceFrequency},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudio, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type oscVoice struct {
	synth *core.Synth // top-level synth
	adsr  core.Module // adsr envelope
	osc   core.Module // oscillator
}

// NewOsc returns an oscillator voice module.
func NewOsc(s *core.Synth, osc core.Module) core.Module {
	log.Info.Printf("")

	adsr := env.NewADSREnv(s)

	// oscillator defaults
	core.SendEventFloat(osc, "duty", 0.1)
	core.SendEventFloat(osc, "attenuation", 1.0)
	core.SendEventFloat(osc, "slope", 0.5)
	// adsr defaults
	core.SendEventFloat(adsr, "attack", 0.1)
	core.SendEventFloat(adsr, "decay", 0.5)
	core.SendEventFloat(adsr, "sustain", 0.05)
	core.SendEventFloat(adsr, "release", 1)

	return &oscVoice{
		synth: s,
		adsr:  adsr,
		osc:   osc,
	}
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
	core.SendEvent(m.adsr, "gate", e)
}

func oscVoiceFrequency(cm core.Module, e *core.Event) {
	m := cm.(*oscVoice)
	core.SendEvent(m.osc, "frequency", e)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *oscVoice) Process(buf ...*core.Buf) {
	out := buf[0]
	// generate wave
	m.osc.Process(out)
	// generate envelope
	var env core.Buf
	m.adsr.Process(&env)
	// apply envelope
	out.Mul(&env)
}

// Active returns true if the module has non-zero output.
func (m *oscVoice) Active() bool {
	return m.adsr.Active()
}

//-----------------------------------------------------------------------------
