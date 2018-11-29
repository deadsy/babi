//-----------------------------------------------------------------------------
/*

Karplus Strong Voice

This provides some controls and defaults for a generic karplus strong oscillator.

*/
//-----------------------------------------------------------------------------

package voice

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var ksVoiceInfo = core.ModuleInfo{
	Name: "ksVoice",
	In: []core.PortInfo{
		{"gate", "oscillator gate, attack(>0) or mute(=0)", core.PortTypeFloat, ksVoiceGate},
		{"frequency", "frequency (Hz)", core.PortTypeFloat, ksVoiceFrequency},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *ksVoice) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type ksVoice struct {
	info core.ModuleInfo // module info
	ks   core.Module     // karplus strong oscillator
}

// NewKarplusStrong returns an karplus strong voice module.
func NewKarplusStrong(s *core.Synth) core.Module {
	log.Info.Printf("new voice")

	ks := osc.NewKarplusStrong(s)
	// ks default
	core.SendEventFloat(ks, "attenuation", 1.0)

	m := &ksVoice{
		info: ksVoiceInfo,
		ks:   ks,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *ksVoice) Child() []core.Module {
	return []core.Module{m.ks}
}

// Stop performs any cleanup of a module.
func (m *ksVoice) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func ksVoiceGate(cm core.Module, e *core.Event) {
	m := cm.(*ksVoice)
	core.SendEvent(m.ks, "gate", e)
}

func ksVoiceFrequency(cm core.Module, e *core.Event) {
	m := cm.(*ksVoice)
	core.SendEvent(m.ks, "frequency", e)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *ksVoice) Process(buf ...*core.Buf) {
	out := buf[0]
	m.ks.Process(out)
}

// Active returns true if the module has non-zero output.
func (m *ksVoice) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
