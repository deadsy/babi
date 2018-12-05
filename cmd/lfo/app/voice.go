//-----------------------------------------------------------------------------
/*

LFO test voice

*/
//-----------------------------------------------------------------------------

package app

import (
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

type voiceApp struct {
	info     core.ModuleInfo // module info
	wavOsc   core.Module     // wave oscillator
	ampEnv   core.Module     // amplitude envelope generator
	velocity float32         // note velocity
}

// NewVoice returns a Goom voice.
func NewVoice(s *core.Synth) core.Module {
	log.Info.Printf("")

	wavOsc := osc.NewGoom(s)
	ampEnv := env.NewADSR(s)

	m := &voiceApp{
		info:   voiceAppInfo,
		wavOsc: wavOsc,
		ampEnv: ampEnv,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *voiceApp) Child() []core.Module {
	return []core.Module{m.ampEnv, m.wavOsc}
}

// Stop performs any cleanup of a module.
func (m *voiceApp) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func voiceAppNote(cm core.Module, e *core.Event) {
	m := cm.(*voiceApp)
	note := e.GetEventFloat().Val
	// set the wave oscillator frequency
	core.EventInFloat(m.wavOsc, "frequency", core.MIDIToFrequency(note))
}

func voiceAppGate(cm core.Module, e *core.Event) {
	m := cm.(*voiceApp)
	gate := e.GetEventFloat().Val
	log.Info.Printf("gate %f", gate)
	if gate > 0 {
		// gate all of the envelopes
		core.EventInFloat(m.ampEnv, "gate", gate)
		// record the note velocity
		m.velocity = gate
	} else {
		// release all of the envelopes
		core.EventInFloat(m.ampEnv, "gate", 0)
	}
}

func voiceAppMidiIn(cm core.Module, e *core.Event) {
	m := cm.(*voiceApp)
	me := e.GetEventMIDI()
	if me != nil {
		if me.GetType() == core.EventMIDIControlChange {

			//val := me.GetCcInt()
			fval := me.GetCcFloat()

			switch me.GetCcNum() {

			case midiWaveDutyCC: // wave oscillator duty cycle
				core.EventInFloat(m.wavOsc, "duty", fval)

			case midiWaveSlopeCC: // wave oscillator slope
				core.EventInFloat(m.wavOsc, "slope", fval)

			case midiAmpAttackCC: // amplitude attack (secs)
				core.EventInFloat(m.ampEnv, "attack", core.MapLin(fval, 0.01, 0.4))

			case midiAmpDecayCC: // amplitude decay (secs)
				core.EventInFloat(m.ampEnv, "decay", core.MapLin(fval, 0.01, 2.0))

			case midiAmpSustainCC: // amplitude sustain (0..1)
				core.EventInFloat(m.ampEnv, "sustain", fval)

			case midiAmpReleaseCC: // amplitude release (secs)
				core.EventInFloat(m.ampEnv, "release", core.MapLin(fval, 0.02, 2.0))

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

	// generate wave
	m.wavOsc.Process(out)

	// generate envelope
	var env core.Buf
	m.ampEnv.Process(&env)

	// apply the envelope
	out.Mul(&env)
}

// Active returns true if the module has non-zero output.
func (m *voiceApp) Active() bool {
	return m.ampEnv.Active()
}

//-----------------------------------------------------------------------------
