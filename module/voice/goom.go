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
	gvPortNull = iota
	gvPortNote
	gvPortGate
	gvPortOscillatorMode
	gvPortFrequencyMode
	gvPortAmplitudeAttack
	gvPortAmplitudeDecay
	gvPortAmplitudeSustain
	gvPortAmplitudeRelease
	gvPortWaveDuty
	gvPortWaveSlope
	gvPortModulationAttack
	gvPortModulationDecay
	gvPortModulationDuty
	gvPortModulationSlope
	gvPortModulationTuning
	gvPortModulationLevel
	gvPortFilterAttack
	gvPortFilterDecay
	gvPortFilterSustain
	gvPortFilterRelease
	gvPortFilterSensitivity
	gvPortFilterCutoff
	gvPortFilterResonance
)

// Info returns the module information.
func (m *goomVoice) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "goomVoice",
		In: []core.PortInfo{
			// overall control
			{"note", "note value (midi)", core.PortTypeFloat, gvPortNote},
			{"gate", "voice gate, attack(>0) or release(=0)", core.PortTypeFloat, gvPortGate},
			{"omode", "oscillator combine mode (0,1,2)", core.PortTypeInt, gvPortOscillatorMode},
			{"fmode", "frequency mode (0,1,2)", core.PortTypeInt, gvPortFrequencyMode},
			// wave envelope
			{"amp_attack", "amplitude attack time (secs)", core.PortTypeFloat, gvPortAmplitudeAttack},
			{"amp_decay", "amplitude decay time (secs)", core.PortTypeFloat, gvPortAmplitudeDecay},
			{"amp_sustain", "amplitude sustain level 0..1", core.PortTypeFloat, gvPortAmplitudeSustain},
			{"amp_release", "amplitude release time (secs)", core.PortTypeFloat, gvPortAmplitudeRelease},
			// wave oscillator
			{"wav_duty", "wave duty cycle (0..1)", core.PortTypeFloat, gvPortWaveDuty},
			{"wav_slope", "wave slope (0..1)", core.PortTypeFloat, gvPortWaveSlope},
			// modulation envelope
			{"mod_attack", "modulation attack time (secs)", core.PortTypeFloat, gvPortModulationAttack},
			{"mod_decay", "modulation decay time (secs)", core.PortTypeFloat, gvPortModulationDecay},
			// modulation oscillator
			{"mod_duty", "modulation duty cycle (0..1)", core.PortTypeFloat, gvPortModulationDuty},
			{"mod_slope", "modulation slope (0..1)", core.PortTypeFloat, gvPortModulationSlope},
			// modulation control
			{"mod_tuning", "modulation tuning (0..1)", core.PortTypeFloat, gvPortModulationTuning},
			{"mod_level", "modulation level (0..1)", core.PortTypeFloat, gvPortModulationLevel},
			// filter envelope
			{"flt_attack", "filter attack time (secs)", core.PortTypeFloat, gvPortFilterAttack},
			{"flt_decay", "filter decay time (secs)", core.PortTypeFloat, gvPortFilterDecay},
			{"flt_sustain", "filter sustain level 0..1", core.PortTypeFloat, gvPortFilterSustain},
			{"flt_release", "filter release time (secs)", core.PortTypeFloat, gvPortFilterRelease},
			// filter control
			{"flt_sensitivity", "low pass filter sensitivity", core.PortTypeFloat, gvPortFilterSensitivity},
			{"flt_cutoff", "low pass filter cutoff frequency (Hz)", core.PortTypeFloat, gvPortFilterCutoff},
			{"flt_resonance", "low pass filter resonance (0..1)", core.PortTypeFloat, gvPortFilterResonance},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

type oModeType uint

const (
	oMode0 oModeType = iota
	oMode1
	oMode2
	oModeMax // must be last
)

type fModeType uint

const (
	fModeLow fModeType = iota
	fModeHigh
	fModeNote
	fModeMax // must be last
)

type goomVoice struct {
	synth  *core.Synth // top-level synth
	wavEnv core.Module // wave envelope generator
	wavOsc core.Module // wave oscillator
	modEnv core.Module // modulation envelope generator
	modOsc core.Module // modulation oscillator
	fltEnv core.Module // filter envelope generator
	lpf    core.Module // low pass filter
	oMode  oModeType   // oscillator mode
	fMode  fModeType   // frequency mode
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

func goomPortOscillatorMode(m *goomVoice, e *core.Event) {
	val := e.GetEventInt().Val
	if !core.InEnum(val, int(oModeMax)) {
		log.Info.Printf("bad value for oscillator mode %d", val)
		return
	}
	m.oMode = oModeType(val)
}

func goomPortFrequencyrMode(m *goomVoice, e *core.Event) {
	val := e.GetEventInt().Val
	if !core.InEnum(val, int(fModeMax)) {
		log.Info.Printf("bad value for frequency mode %d", val)
	}
	m.fMode = fModeType(val)
}

func goomPortNote(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortGate(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortAmplitudeAttack(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortAmplitudeDecay(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortAmplitudeSustain(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortAmplitudeRelease(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortWaveDuty(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortWaveSlope(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortModulationAttack(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortModulationDecay(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortModulationDuty(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortModulationSlope(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortModulationTuning(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortModulationLevel(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortFilterAttack(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortFilterDecay(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortFilterSustain(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortFilterRelease(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortFilterSensitivity(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortFilterCutoff(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

func goomPortFilterResonance(m *goomVoice, e *core.Event) {
	val := e.GetEventFloat().Val
	_ = val
}

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
