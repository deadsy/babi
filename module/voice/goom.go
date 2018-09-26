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

// Info returns the module information.
func (m *goomVoice) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "goomVoice",
		In: []core.PortInfo{
			// overall control
			{"note", "note value (midi)", core.PortTypeFloat, goomPortNote},
			{"gate", "voice gate, attack(>0) or release(=0)", core.PortTypeFloat, goomPortGate},
			{"omode", "oscillator combine mode (0,1,2)", core.PortTypeInt, goomPortOscillatorMode},
			{"fmode", "frequency mode (0,1,2)", core.PortTypeInt, goomPortFrequencyMode},
			// amplitude envelope
			{"amp_attack", "amplitude attack time (secs)", core.PortTypeFloat, goomPortAmplitudeAttack},
			{"amp_decay", "amplitude decay time (secs)", core.PortTypeFloat, goomPort},
			{"amp_sustain", "amplitude sustain level 0..1", core.PortTypeFloat, goomPort},
			{"amp_release", "amplitude release time (secs)", core.PortTypeFloat, goomPort0},
			// wave oscillator
			{"wav_duty", "wave duty cycle (0..1)", core.PortTypeFloat, goomPort0},
			{"wav_slope", "wave slope (0..1)", core.PortTypeFloat, goomPort0},
			// modulation envelope
			{"mod_attack", "modulation attack time (secs)", core.PortTypeFloat, goomPort0},
			{"mod_decay", "modulation decay time (secs)", core.PortTypeFloat, goomPort0},
			// modulation oscillator
			{"mod_duty", "modulation duty cycle (0..1)", core.PortTypeFloat, goomPort0},
			{"mod_slope", "modulation slope (0..1)", core.PortTypeFloat, goomPort0},
			// modulation control
			{"mod_tuning", "modulation tuning (0..1)", core.PortTypeFloat, goomPort0},
			{"mod_level", "modulation level (0..1)", core.PortTypeFloat, goomPort0},
			// filter envelope
			{"flt_attack", "filter attack time (secs)", core.PortTypeFloat, goomPort0},
			{"flt_decay", "filter decay time (secs)", core.PortTypeFloat, goomPort0},
			{"flt_sustain", "filter sustain level 0..1", core.PortTypeFloat, goomPort0},
			{"flt_release", "filter release time (secs)", core.PortTypeFloat, goomPort0},
			// filter control
			{"flt_sensitivity", "low pass filter sensitivity", core.PortTypeFloat, goomPort0},
			{"flt_cutoff", "low pass filter cutoff frequency (Hz)", core.PortTypeFloat, goomPort0},
			{"flt_resonance", "low pass filter resonance (0..1)", core.PortTypeFloat, goomPort0},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, nil},
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
	synth     *core.Synth // top-level synth
	ampEnv    core.Module // amplitude envelope generator
	wavOsc    core.Module // wave oscillator
	modEnv    core.Module // modulation envelope generator
	modOsc    core.Module // modulation oscillator
	fltEnv    core.Module // filter envelope generator
	lpf       core.Module // low pass filter
	oMode     oModeType   // oscillator mode
	fMode     fModeType   // frequency mode
	modTuning float32     // modulation tuning
	modLevel  float32     // modulation level
}

// NewGoomVoice returns a Goom voice.
func NewGoomVoice(s *core.Synth) core.Module {
	log.Info.Printf("")

	ampEnv := env.NewADSREnv(s)
	wavOsc := osc.NewGoomOsc(s)
	modEnv := env.NewADSREnv(s)
	modOsc := osc.NewGoomOsc(s)
	fltEnv := env.NewADSREnv(s)
	lpf := filter.NewSVFilterTrapezoidal(s)

	return &goomVoice{
		synth:  s,
		ampEnv: ampEnv,
		wavOsc: wavOsc,
		modEnv: modEnv,
		modOsc: modOsc,
		fltEnv: fltEnv,
		lpf:    lpf,
	}
}

// Child returns the child modules of this module.
func (m *goomVoice) Child() []core.Module {
	return []core.Module{m.ampEnv, m.wavOsc, m.modEnv, m.modOsc, m.fltEnv, m.lpf}
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
	log.Info.Printf("set oscillator mode %d", val)
	m.oMode = oModeType(val)
}

func goomPortFrequencyrMode(m *goomVoice, e *core.Event) {
	val := e.GetEventInt().Val
	if !core.InEnum(val, int(fModeMax)) {
		log.Info.Printf("bad value for frequency mode %d", val)
		return
	}
	log.Info.Printf("set frequency mode %d", val)
	m.fMode = fModeType(val)
}

func goomPortNote(m *goomVoice, e *core.Event) {
	// TODO
}

func goomPortGate(m *goomVoice, e *core.Event) {
	// TODO
}

func goomPortAmplitudeAttack(m *goomVoice, e *core.Event) {
	core.SendEvent(m.ampEnv, "attack", e)
}

func goomPortAmplitudeDecay(m *goomVoice, e *core.Event) {
	core.SendEvent(m.ampEnv, "decay", e)
}

func goomPortAmplitudeSustain(m *goomVoice, e *core.Event) {
	core.SendEvent(m.ampEnv, "sustain", e)
}

func goomPortAmplitudeRelease(m *goomVoice, e *core.Event) {
	core.SendEvent(m.ampEnv, "release", e)
}

func goomPortWaveDuty(m *goomVoice, e *core.Event) {
	core.SendEvent(m.wavOsc, "duty", e)
}

func goomPortWaveSlope(m *goomVoice, e *core.Event) {
	core.SendEvent(m.wavOsc, "slope", e)
}

func goomPortModulationAttack(m *goomVoice, e *core.Event) {
	core.SendEvent(m.modEnv, "attack", e)
}

func goomPortModulationDecay(m *goomVoice, e *core.Event) {
	core.SendEvent(m.modEnv, "decay", e)
}

func goomPortModulationDuty(m *goomVoice, e *core.Event) {
	core.SendEvent(m.modOsc, "duty", e)
}

func goomPortModulationSlope(m *goomVoice, e *core.Event) {
	core.SendEvent(m.modOsc, "slope", e)
}

func goomPortModulationTuning(m *goomVoice, e *core.Event) {
	m.modTuning = core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set modulation tuning %f", m.modTuning)
}

func goomPortModulationLevel(m *goomVoice, e *core.Event) {
	m.modLevel = core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set modulation level %f", m.modLevel)
}

func goomPortFilterAttack(m *goomVoice, e *core.Event) {
	core.SendEvent(m.fltEnv, "attack", e)
}

func goomPortFilterDecay(m *goomVoice, e *core.Event) {
	core.SendEvent(m.fltEnv, "decay", e)
}

func goomPortFilterSustain(m *goomVoice, e *core.Event) {
	core.SendEvent(m.fltEnv, "sustain", e)
}

func goomPortFilterRelease(m *goomVoice, e *core.Event) {
	core.SendEvent(m.fltEnv, "release", e)
}

func goomPortFilterSensitivity(m *goomVoice, e *core.Event) {
	// TODO
}

func goomPortFilterCutoff(m *goomVoice, e *core.Event) {
	// TODO
}

func goomPortFilterResonance(m *goomVoice, e *core.Event) {
	core.SendEvent(m.lpf, "resonance", e)
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
