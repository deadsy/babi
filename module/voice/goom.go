//-----------------------------------------------------------------------------
/*

Goom Voice

https://www.quinapalus.com/goom.html

*/
//-----------------------------------------------------------------------------

package module

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/filter"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

const (
	goomPortNull = iota
)

// Info returns the module information.
func (m *goomVoice) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "goomVoice",
		In: []core.PortInfo{
			// overall control
			{"note", "note value (midi)", core.PortTypeFloat, goomPortNull},
			{"gate", "voice gate, attack(>0) or release(=0)", core.PortTypeFloat, goomPortNull},
			{"o_mode", "oscillator combine mode (0,1,2)", core.PortTypeInt, goomPortNull},
			{"f_mode", "frequency mode (0,1,2)", core.PortTypeInt, goomPortNull},
			// wave envelope
			{"amp_attack", "amplitude attack time (secs)", core.PortTypeFloat, goomPortNull},
			{"amp_decay", "amplitude decay time (secs)", core.PortTypeFloat, goomPortNull},
			{"amp_sustain", "amplitude sustain level 0..1", core.PortTypeFloat, goomPortNull},
			{"amp_release", "amplitude release time (secs)", core.PortTypeFloat, goomPortNull},
			// wave oscillator
			{"wav_duty", "wave duty cycle (0..1)", core.PortTypeFloat, goomPortNull},
			{"wav_slope", "wave slope (0..1)", core.PortTypeFloat, goomPortNull},
			// modulation envelope
			{"mod_attack", "modulation attack time (secs)", core.PortTypeFloat, goomPortNull},
			{"mod_decay", "modulation decay time (secs)", core.PortTypeFloat, goomPortNull},
			// modulation oscillator
			{"mod_duty", "modulation duty cycle (0..1)", core.PortTypeFloat, goomPortNull},
			{"mod_slope", "modulation slope (0..1)", core.PortTypeFloat, goomPortNull},
			// modulation control
			{"mod_tuning", "modulation tuning (0..1)", core.PortTypeFloat, goomPortNull},
			{"mod_level", "modulation level (0..1)", core.PortTypeFloat, goomPortNull},
			// filter envelope
			{"flt_attack", "filter attack time (secs)", core.PortTypeFloat, goomPortNull},
			{"flt_decay", "filter decay time (secs)", core.PortTypeFloat, goomPortNull},
			{"flt_sustain", "filter sustain level 0..1", core.PortTypeFloat, goomPortNull},
			{"flt_release", "filter release time (secs)", core.PortTypeFloat, goomPortNull},
			// filter control
			{"flt_sensitivity", "low pass filter sensitivity", core.PortTypeFloat, goomPortNull},
			{"flt_cutoff", "low pass filter cutoff frequency (Hz)", core.PortTypeFloat, goomPortNull},
			{"flt_resonance", "low pass filter resonance (0..1)", core.PortTypeFloat, goomPortNull},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type goomVoice struct {
	synth  *core.Synth // top-level synth
	wavEnv core.Module // wave envelope generator
	wavOsc core.Module // wave oscillator
	modEnv core.Module // modulation envelope generator
	modOsc core.Module // modulation oscillator
	fltEnv core.Module // filter envelope generator
	lpf    core.Module // low pass filter
}

// NewGoomVoice returns a Goom voice.
func NewGoomVoice(s *core.Synth) core.Module {
	log.Info.Printf("")

	wavEnv := env.NewADSREnv(s)
	wavOsc := osc.NewGoomOsc(s)
	modEnv := env.NewADSREnv(s)
	modOsc := osc.NewGoomOsc(s)
	fltEnv := env.NewADSREnv(s)
	lpf := filter.NewSVFilterTrapezoidal(s)

	return &goomVoice{
		synth:  s,
		wavEnv: wavEnv,
		wavOsc: wavOsc,
		modEnv: modEnv,
		modOsc: modOsc,
		fltEnv: fltEnv,
		lpf:    lpf,
	}
}

// Child returns the child modules of this module.
func (m *goomVoice) Child() []core.Module {
	return []core.Module{m.wavEnv, m.wavOsc, m.modEnv, m.modOsc, m.fltEnv, m.lpf}
}

// Stop performs any cleanup of a module.
func (m *goomVoice) Stop() {
}

//-----------------------------------------------------------------------------
// Events

// Event processes a module event.
func (m *goomVoice) Event(e *core.Event) {
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *goomVoice) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *goomVoice) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
