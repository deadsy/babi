//-----------------------------------------------------------------------------
/*

LFO test voice

*/
//-----------------------------------------------------------------------------

package app

import (
	"fmt"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/module/env"
	"github.com/deadsy/babi/module/osc"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

var voiceAppInfo = core.ModuleInfo{
	Name: "voiceApp",
	In: []core.PortInfo{
		// overall control
		{"note", "note value", core.PortTypeFloat, voiceAppNote},
		{"gate", "voice gate, attack(>0) or release(=0)", core.PortTypeFloat, voiceAppGate},
		{"midi", "midi input", core.PortTypeMIDI, voiceAppMidiIn},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *voiceApp) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------
// modulation modes

type modMode int

// modulation modes
const (
	modModeOff modMode = iota // none
	modModeAM                 // amplitude modulation
	modModeFM                 // frequency modulation
	modModePM                 // phase modulation
)

var modModeToString = map[modMode]string{
	modModeOff: "off",
	modModeAM:  "am",
	modModeFM:  "fm",
	modModePM:  "pm",
}

func (m modMode) String() string {
	return modModeToString[m]
}

//-----------------------------------------------------------------------------

type voiceApp struct {
	info     core.ModuleInfo // module info
	lfo      core.Module     // LFO
	wav      core.Module     // wave oscillator
	env      core.Module     // amplitude envelope generator
	mode     modMode         // modulation mode
	note     float32         // midi note (as a float)
	depth    float32         // unscaled lfo depth
	velocity float32         // note velocity
}

// NewVoice returns a LFO test voice.
func NewVoice(s *core.Synth) core.Module {
	log.Info.Printf("")

	m := &voiceApp{
		info: voiceAppInfo,
		lfo:  osc.NewLFO(s),
		wav:  osc.NewGoom(s),
		env:  env.NewADSR(s),
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *voiceApp) Child() []core.Module {
	return []core.Module{m.env, m.wav, m.lfo}
}

// Stop performs any cleanup of a module.
func (m *voiceApp) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

// The LFO depth is set by the CC value but it's scaling depends on the
// note pitch and the modulation type.
func (m *voiceApp) setDepth() {
	var scale float32
	switch m.mode {
	case modModeAM:
		scale = 1.0
	case modModeFM:
		scale = core.MIDIToFrequency(m.note * 0.25)
	case modModePM:
		scale = core.Pi * 0.15
	default:
	}
	core.EventInFloat(m.lfo, "depth", core.MapLin(m.depth, 0, scale))
}

func voiceAppNote(cm core.Module, e *core.Event) {
	m := cm.(*voiceApp)
	note := e.GetEventFloat().Val
	m.note = note
	// set the wave oscillator frequency
	core.EventInFloat(m.wav, "frequency", core.MIDIToFrequency(note))
	// re-set the LFO depth since it is a function of the note
	m.setDepth()
}

func voiceAppGate(cm core.Module, e *core.Event) {
	m := cm.(*voiceApp)
	gate := e.GetEventFloat().Val
	log.Info.Printf("gate %f", gate)
	if gate > 0 {
		// gate the envelope
		core.EventInFloat(m.env, "gate", gate)
		// record the note velocity
		m.velocity = gate
	} else {
		// release the envelope
		core.EventInFloat(m.env, "gate", 0)
	}
}

func voiceAppMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*voiceApp)
	me := e.GetEventMIDI()
	if me != nil {
		if me.GetType() == core.EventMIDIControlChange {
			val := me.GetCcInt()
			fval := me.GetCcFloat()
			switch me.GetCcNum() {
			// wave
			case midiWaveDutyCC: // wave oscillator duty cycle
				core.EventInFloat(m.wav, "duty", fval)
			case midiWaveSlopeCC: // wave oscillator slope
				core.EventInFloat(m.wav, "slope", fval)
			// envelope
			case midiAmpAttackCC: // amplitude attack (secs)
				core.EventInFloat(m.env, "attack", core.MapLin(fval, 0.01, 0.4))
			case midiAmpDecayCC: // amplitude decay (secs)
				core.EventInFloat(m.env, "decay", core.MapLin(fval, 0.01, 2.0))
			case midiAmpSustainCC: // amplitude sustain (0..1)
				core.EventInFloat(m.env, "sustain", fval)
			case midiAmpReleaseCC: // amplitude release (secs)
				core.EventInFloat(m.env, "release", core.MapLin(fval, 0.02, 2.0))
			// lfo
			case midiLfoRateCC: // lfo rate (Hz)
				core.EventInFloat(m.lfo, "rate", core.MapLin(fval, 0.5, 20.0))
			case midiLfoDepthCC: // lfo depth
				m.depth = fval
				m.setDepth()
			case midiLfoShapeCC: // lfo wave shape
				core.EventInInt(m.lfo, "shape", int(val))
			case midiModModeCC: // modulation mode (am/fm/pm)
				m.mode = modMode(val)
				// set the goom oscillator mode
				switch m.mode {
				case modModeOff, modModeAM:
					core.EventInInt(m.wav, "mode", int(osc.GoomModeBasic))
				case modModeFM:
					core.EventInInt(m.wav, "mode", int(osc.GoomModeFM))
				case modModePM:
					core.EventInInt(m.wav, "mode", int(osc.GoomModePM))
				default:
					panic(fmt.Sprintf("bad mode %d", m.mode))
				}
				// re-set the lfo depth since it is a function of mode
				m.setDepth()
			default:
				// ignore
			}
		}
	}
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *voiceApp) Process(buf ...*core.Buf) {
	out := buf[0]

	if m.mode != modModeOff {
		// generate the modulating lfo
		var mod core.Buf
		m.lfo.Process(&mod)
		switch m.mode {
		case modModeAM:
			m.wav.Process(out)
			out.Mul(&mod)
		case modModeFM:
			m.wav.Process(&mod, out)
		case modModePM:
			m.wav.Process(&mod, out)
		default:
			panic(fmt.Sprintf("bad mode %d", m.mode))
		}
	} else {
		m.wav.Process(out)
	}

	// generate envelope
	var env core.Buf
	m.env.Process(&env)

	// apply the envelope
	out.Mul(&env)
}

// Active returns true if the module has non-zero output.
func (m *voiceApp) Active() bool {
	return m.env.Active()
}

//-----------------------------------------------------------------------------
