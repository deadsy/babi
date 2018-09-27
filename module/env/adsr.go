//-----------------------------------------------------------------------------
/*

ADSR Envelope Module

*/
//-----------------------------------------------------------------------------

package env

import (
	"fmt"
	"math"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *adsrModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "adsr",
		In: []core.PortInfo{
			{"gate", "envelope gate, attack(>0) or release(=0)", core.PortTypeFloat, adsrPortGate},
			{"attack", "attack time (secs)", core.PortTypeFloat, adsrPortAttack},
			{"decay", "decay time (secs)", core.PortTypeFloat, adsrPortDecay},
			{"sustain", "sustain level 0..1", core.PortTypeFloat, adsrPortSustain},
			{"release", "release time (secs)", core.PortTypeFloat, adsrPortRelease},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudioBuffer, nil},
		},
	}
}

//-----------------------------------------------------------------------------

// We can't reach the target level with the asymptotic rise/fall of exponentials.
// We will change state when we are within adsrEpsilon of the target level.
const adsrEpsilon = 0.001

// Return a k value to give the exponential rise/fall in the required time.
func getK(t float32, rate int) float32 {
	if t <= 0 {
		return 1.0
	}
	return float32(1.0 - math.Exp(math.Log(adsrEpsilon)/(float64(t)*float64(rate))))
}

//-----------------------------------------------------------------------------

type adsrState int

const (
	stateIdle adsrState = iota // initial state
	stateAttack
	stateDecay
	stateSustain
	stateRelease
)

type adsrModule struct {
	synth    *core.Synth // top-level synth
	state    adsrState   // envelope state
	s        float32     // sustain level
	ka       float32     // attack constant
	kd       float32     // decay constant
	kr       float32     // release constant
	dTrigger float32     // attack->decay trigger level
	sTrigger float32     // decay->sustain trigger level
	iTrigger float32     // release->idle trigger level
	val      float32     // output value
}

// NewADSREnv returns an Attack/Decay/Sustain/Release envelope module.
func NewADSREnv(s *core.Synth) core.Module {
	log.Info.Printf("")
	return &adsrModule{
		synth: s,
	}
}

// Return the child modules.
func (m *adsrModule) Child() []core.Module {
	return nil
}

// Stop stops and performs any cleanup of a module.
func (m *adsrModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Port Events

func adsrPortGate(cm core.Module, e *core.Event) {
	m := cm.(*adsrModule)
	gate := e.GetEventFloat().Val
	log.Info.Printf("gate %f", gate)
	if gate != 0 {
		// enter the attack segment
		m.state = stateAttack
	} else {
		// enter the release segment
		if m.state != stateIdle {
			if m.kr == 1 {
				// no release - goto idle
				m.val = 0
				m.state = stateIdle
			} else {
				m.state = stateRelease
			}
		}
	}
}

func adsrPortAttack(cm core.Module, e *core.Event) {
	m := cm.(*adsrModule)
	attack := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set attack time %f secs", attack)
	m.ka = getK(attack, core.AudioSampleFrequency)
}

func adsrPortDecay(cm core.Module, e *core.Event) {
	m := cm.(*adsrModule)
	decay := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set decay time %f secs", decay)
	m.kd = getK(decay, core.AudioSampleFrequency)
}

func adsrPortSustain(cm core.Module, e *core.Event) {
	m := cm.(*adsrModule)
	sustain := core.Clamp(e.GetEventFloat().Val, 0, 1)
	log.Info.Printf("set sustain level %f", sustain)
	m.s = sustain
	m.dTrigger = 1.0 - adsrEpsilon
	m.sTrigger = sustain + (1.0-sustain)*adsrEpsilon
	m.iTrigger = sustain * adsrEpsilon
}

func adsrPortRelease(cm core.Module, e *core.Event) {
	m := cm.(*adsrModule)
	release := core.ClampLo(e.GetEventFloat().Val, 0)
	log.Info.Printf("set release time %f secs", release)
	m.kr = getK(release, core.AudioSampleFrequency)
}

//-----------------------------------------------------------------------------

// Process runs the module DSP.
func (m *adsrModule) Process(buf ...*core.Buf) {
	out := buf[0]
	for i := 0; i < len(out); i++ {
		switch m.state {
		case stateIdle:
			// idle - do nothing
		case stateAttack:
			// attack until 1.0 level
			if m.val < m.dTrigger {
				m.val += m.ka * (1.0 - m.val)
			} else {
				// goto decay state
				m.val = 1
				m.state = stateDecay
			}
		case stateDecay:
			// decay until sustain level
			if m.val > m.sTrigger {
				m.val += m.kd * (m.s - m.val)
			} else {
				if m.s != 0 {
					// goto sustain state
					m.val = m.s
					m.state = stateSustain
				} else {
					// no sustain, goto idle state
					m.val = 0
					m.state = stateIdle
				}
			}
		case stateSustain:
			// sustain - do nothing
		case stateRelease:
			// release until idle level
			if m.val > m.iTrigger {
				m.val += m.kr * (0.0 - m.val)
			} else {
				// goto idle state
				m.val = 0
				m.state = stateIdle
			}
		default:
			panic(fmt.Sprintf("bad adsr state %d", m.state))
		}
		out[i] = m.val
	}
}

// Active return true if the module has non-zero output.
func (m *adsrModule) Active() bool {
	return m.state != stateIdle
}

//-----------------------------------------------------------------------------
