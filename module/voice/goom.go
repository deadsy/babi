//-----------------------------------------------------------------------------
/*

Goom Voice

https://www.quinapalus.com/goom.html

*/
//-----------------------------------------------------------------------------

package voice

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
			{"amp_decay", "amplitude decay time (secs)", core.PortTypeFloat, goomPortAmplitudeDecay},
			{"amp_sustain", "amplitude sustain level 0..1", core.PortTypeFloat, goomPortAmplitudeSustain},
			{"amp_release", "amplitude release time (secs)", core.PortTypeFloat, goomPortAmplitudeRelease},
			// wave oscillator
			{"wav_duty", "wave duty cycle (0..1)", core.PortTypeFloat, goomPortWaveDuty},
			{"wav_slope", "wave slope (0..1)", core.PortTypeFloat, goomPortWaveSlope},
			// modulation envelope
			{"mod_attack", "modulation attack time (secs)", core.PortTypeFloat, goomPortModulationAttack},
			{"mod_decay", "modulation decay time (secs)", core.PortTypeFloat, goomPortModulationDecay},
			// modulation oscillator
			{"mod_duty", "modulation duty cycle (0..1)", core.PortTypeFloat, goomPortModulationDuty},
			{"mod_slope", "modulation slope (0..1)", core.PortTypeFloat, goomPortModulationSlope},
			// modulation control
			{"mod_tuning", "modulation tuning (0..1)", core.PortTypeFloat, goomPortModulationTuning},
			{"mod_level", "modulation level (0..1)", core.PortTypeFloat, goomPortModulationLevel},
			// filter envelope
			{"flt_attack", "filter attack time (secs)", core.PortTypeFloat, goomPortFilterAttack},
			{"flt_decay", "filter decay time (secs)", core.PortTypeFloat, goomPortFilterDecay},
			{"flt_sustain", "filter sustain level 0..1", core.PortTypeFloat, goomPortFilterSustain},
			{"flt_release", "filter release time (secs)", core.PortTypeFloat, goomPortFilterRelease},
			// filter control
			{"flt_sensitivity", "low pass filter sensitivity", core.PortTypeFloat, goomPortFilterSensitivity},
			{"flt_cutoff", "low pass filter cutoff frequency (Hz)", core.PortTypeFloat, goomPortFilterCutoff},
			{"flt_resonance", "low pass filter resonance (0..1)", core.PortTypeFloat, goomPortFilterResonance},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudio, nil},
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
	fModeNote fModeType = iota
	fModeHigh
	fModeLow
	fModeMax // must be last
)

type goomVoice struct {
	synth          *core.Synth // top-level synth
	ampEnv         core.Module // amplitude envelope generator
	wavOsc         core.Module // wave oscillator
	modEnv         core.Module // modulation envelope generator
	modOsc         core.Module // modulation oscillator
	fltEnv         core.Module // filter envelope generator
	lpf            core.Module // low pass filter
	oMode          oModeType   // oscillator mode
	fMode          fModeType   // frequency mode
	modTuning      float32     // modulation tuning
	modLevel       float32     // modulation level
	fltSensitivity float32     // filter sensitivity
	fltCutoff      float32     // filter cutoff
	velocity       float32     // note velocity
}

// NewGoom returns a Goom voice.
func NewGoom(s *core.Synth) core.Module {
	log.Info.Printf("")

	ampEnv := env.NewADSREnv(s)
	wavOsc := osc.NewGoom(s)
	modEnv := env.NewADSREnv(s)
	modOsc := osc.NewGoom(s)
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
// Port Events

func goomPortOscillatorMode(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	val := e.GetEventInt().Val
	if !core.InEnum(val, int(oModeMax)) {
		log.Info.Printf("bad value for oscillator mode %d", val)
		return
	}
	log.Info.Printf("set oscillator mode %d", val)
	m.oMode = oModeType(val)
}

func goomPortFrequencyMode(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	val := e.GetEventInt().Val
	if !core.InEnum(val, int(fModeMax)) {
		log.Info.Printf("bad value for frequency mode %d", val)
		return
	}
	log.Info.Printf("set frequency mode %d", val)
	m.fMode = fModeType(val)
}

func goomPortNote(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	note := e.GetEventFloat().Val
	// set the wave oscillator frequency
	core.SendEventFloat(m.wavOsc, "frequency", core.MIDIToFrequency(note))
	// set the modulation oscillator frequency
	switch m.fMode {
	case fModeLow:
		note = 10
	case fModeHigh:
		note = 100
	}
	note += m.modTuning * 2 // +/- 2 semitones
	core.SendEventFloat(m.modOsc, "frequency", core.MIDIToFrequency(note))
}

func goomPortGate(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	gate := e.GetEventFloat().Val
	log.Info.Printf("gate %f", gate)
	if gate > 0 {
		// gate all of the envelopes
		core.SendEventFloat(m.ampEnv, "gate", gate)
		core.SendEventFloat(m.modEnv, "gate", gate)
		core.SendEventFloat(m.fltEnv, "gate", gate)
		// record the note velocity
		m.velocity = gate
	} else {
		// release all of the envelopes
		core.SendEventFloat(m.ampEnv, "gate", 0)
		core.SendEventFloat(m.modEnv, "gate", 0)
		core.SendEventFloat(m.fltEnv, "gate", 0)
	}
}

func goomPortAmplitudeAttack(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.ampEnv, "attack", e)
}

func goomPortAmplitudeDecay(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.ampEnv, "decay", e)
}

func goomPortAmplitudeSustain(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.ampEnv, "sustain", e)
}

func goomPortAmplitudeRelease(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.ampEnv, "release", e)
}

func goomPortWaveDuty(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.wavOsc, "duty", e)
}

func goomPortWaveSlope(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.wavOsc, "slope", e)
}

func goomPortModulationAttack(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.modEnv, "attack", e)
}

func goomPortModulationDecay(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.modEnv, "decay", e)
}

func goomPortModulationDuty(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.modOsc, "duty", e)
}

func goomPortModulationSlope(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.modOsc, "slope", e)
}

func goomPortModulationTuning(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	tune := core.Clamp(e.GetEventFloat().Val, 0, 1)
	tune = core.Map(tune, -1, 1)
	log.Info.Printf("set modulation tuning %f", tune)
	m.modTuning = tune
}

func goomPortModulationLevel(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	m.modLevel = core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set modulation level %f", m.modLevel)
}

func goomPortFilterAttack(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.fltEnv, "attack", e)
}

func goomPortFilterDecay(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.fltEnv, "decay", e)
}

func goomPortFilterSustain(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.fltEnv, "sustain", e)
}

func goomPortFilterRelease(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.fltEnv, "release", e)
}

func goomPortFilterSensitivity(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	sensitivity := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set filter sensitivity %f", sensitivity)
	m.fltSensitivity = sensitivity
}

func goomPortFilterCutoff(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	cutoff := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set filter cutoff %f", cutoff)
	m.fltCutoff = cutoff
}

func goomPortFilterResonance(cm core.Module, e *core.Event) {
	m := cm.(*goomVoice)
	core.SendEvent(m.lpf, "resonance", e)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *goomVoice) Process(buf ...*core.Buf) {
}

// Active returns true if the module has non-zero output.
func (m *goomVoice) Active() bool {
	return m.ampEnv.Active()
}

//-----------------------------------------------------------------------------
