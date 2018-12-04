//-----------------------------------------------------------------------------
/*

Goom Voice

https://www.quinapalus.com/goom.html

*/
//-----------------------------------------------------------------------------

package goom

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/filter"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var voiceGoomInfo = core.ModuleInfo{
	Name: "voiceGoom",
	In: []core.PortInfo{
		// overall control
		{"note", "note value", core.PortTypeFloat, voiceGoomNote},
		{"gate", "voice gate, attack(>0) or release(=0)", core.PortTypeFloat, voiceGoomGate},
		{"midi", "midi input", core.PortTypeMIDI, voiceGoomMidiIn},

		/*
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
		*/
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *voiceGoom) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type oModeType uint

const (
	oMode0 oModeType = iota
	oMode1
	oMode2
)

type fModeType uint

const (
	fModeNote fModeType = iota
	fModeHigh
	fModeLow
)

type voiceGoom struct {
	info   core.ModuleInfo // module info
	wavOsc core.Module     // wave oscillator
	ampEnv core.Module     // amplitude envelope generator
	oMode  oModeType       // oscillator mode
	fMode  fModeType       // frequency mode

	modEnv         core.Module // modulation envelope generator
	fltEnv         core.Module // filter envelope generator
	modOsc         core.Module // modulation oscillator
	lpf            core.Module // low pass filter
	modTuning      float32     // modulation tuning
	modLevel       float32     // modulation level
	fltSensitivity float32     // filter sensitivity
	fltCutoff      float32     // filter cutoff
	velocity       float32     // note velocity
}

// NewVoice returns a Goom voice.
func NewVoice(s *core.Synth) core.Module {
	log.Info.Printf("")

	ampEnv := env.NewADSR(s)
	wavOsc := osc.NewGoom(s)
	modEnv := env.NewADSR(s)
	modOsc := osc.NewGoom(s)
	fltEnv := env.NewADSR(s)
	lpf := filter.NewSVFilterTrapezoidal(s)

	m := &voiceGoom{
		info:   voiceGoomInfo,
		ampEnv: ampEnv,
		wavOsc: wavOsc,
		modEnv: modEnv,
		modOsc: modOsc,
		fltEnv: fltEnv,
		lpf:    lpf,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *voiceGoom) Child() []core.Module {
	return []core.Module{m.ampEnv, m.wavOsc, m.modEnv, m.modOsc, m.fltEnv, m.lpf}
}

// Stop performs any cleanup of a module.
func (m *voiceGoom) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func voiceGoomNote(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	note := e.GetEventFloat().Val
	// set the wave oscillator frequency
	core.EventInFloat(m.wavOsc, "frequency", core.MIDIToFrequency(note))

	/*
		  // set the modulation oscillator frequency
			switch m.fMode {
			case fModeLow:
				note = 10
			case fModeHigh:
				note = 100
			}
			note += m.modTuning * 2 // +/- 2 semitones
			core.EventInFloat(m.modOsc, "frequency", core.MIDIToFrequency(note))
	*/
}

func voiceGoomGate(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	gate := e.GetEventFloat().Val
	log.Info.Printf("gate %f", gate)
	if gate > 0 {
		// gate all of the envelopes
		core.EventInFloat(m.ampEnv, "gate", gate)
		core.EventInFloat(m.modEnv, "gate", gate)
		core.EventInFloat(m.fltEnv, "gate", gate)
		// record the note velocity
		m.velocity = gate
	} else {
		// release all of the envelopes
		core.EventInFloat(m.ampEnv, "gate", 0)
		core.EventInFloat(m.modEnv, "gate", 0)
		core.EventInFloat(m.fltEnv, "gate", 0)
	}
}

func voiceGoomMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	me := e.GetEventMIDI()
	if me != nil {
		if me.GetType() == core.EventMIDIControlChange {

			val := me.GetCcInt()
			fval := me.GetCcFloat()

			switch me.GetCcNum() {

			case midiWaveDutyCC:
				// wave oscillator duty cycle
				core.EventInFloat(m.wavOsc, "duty", fval)

			case midiWaveSlopeCC:
				// wave oscillator duty slope
				core.EventInFloat(m.wavOsc, "slope", fval)

			case midiAmpAttackCC:
				// amplitude attack (secs)
				core.EventInFloat(m.ampEnv, "attack", core.MapLin(fval, 0.01, 0.4))

			case midiAmpDecayCC:
				// amplitude decay (secs)
				core.EventInFloat(m.ampEnv, "decay", core.MapLin(fval, 0.01, 2.0))

			case midiAmpSustainCC:
				// amplitude sustain (0..1)
				core.EventInFloat(m.ampEnv, "sustain", fval)

			case midiAmpReleaseCC:
				// amplitude release (secs)
				core.EventInFloat(m.ampEnv, "release", core.MapLin(fval, 0.02, 2.0))

			case midiOscillatorModeCC:
				// oscillator combine mode
				log.Info.Printf("set oscillator mode %d", val)
				m.oMode = oModeType(val)

			case midiFrequencyModeCC:
				// frequency mode
				log.Info.Printf("set frequency mode %d", val)
				m.fMode = fModeType(val)

			default:
				// ignore
			}
		}
	}
}

/*

func goomPortModulationAttack(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.modEnv, "attack", e)
}

func goomPortModulationDecay(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.modEnv, "decay", e)
}

func goomPortModulationDuty(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.modOsc, "duty", e)
}

func goomPortModulationSlope(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.modOsc, "slope", e)
}

func goomPortModulationTuning(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	tune := core.Clamp(e.GetEventFloat().Val, 0, 1)
	tune = core.Map(tune, -1, 1)
	log.Info.Printf("set modulation tuning %f", tune)
	m.modTuning = tune
}

func goomPortModulationLevel(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	m.modLevel = core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set modulation level %f", m.modLevel)
}

func goomPortFilterAttack(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.fltEnv, "attack", e)
}

func goomPortFilterDecay(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.fltEnv, "decay", e)
}

func goomPortFilterSustain(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.fltEnv, "sustain", e)
}

func goomPortFilterRelease(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.fltEnv, "release", e)
}

func goomPortFilterSensitivity(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	sensitivity := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set filter sensitivity %f", sensitivity)
	m.fltSensitivity = sensitivity
}

func goomPortFilterCutoff(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	cutoff := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set filter cutoff %f", cutoff)
	m.fltCutoff = cutoff
}

func goomPortFilterResonance(cm core.Module, e *core.Event) {
	m := cm.(*voiceGoom)
	core.EventIn(m.lpf, "resonance", e)
}

*/

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *voiceGoom) Process(buf ...*core.Buf) {
	out := buf[0]
	// generate wave
	m.wavOsc.Process(out)

	// generate envelope
	var env core.Buf
	m.ampEnv.Process(&env)

	// apply envelope
	out.Mul(&env)
}

// Active returns true if the module has non-zero output.
func (m *voiceGoom) Active() bool {
	return m.ampEnv.Active()
}

//-----------------------------------------------------------------------------
