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
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *adsrModule) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "adsr",
		In: []core.PortInfo{
			{"gate", "envelope gate, attack(>0) or release(=0)", core.PortType_EventFloat, 0},
			{"attack", "attack time (secs)", core.PortType_EventFloat, 0},
			{"decay", "decay time (secs)", core.PortType_EventFloat, 0},
			{"sustain", "sustain level 0..1", core.PortType_EventFloat, 0},
			{"release", "release time (secs)", core.PortType_EventFloat, 0},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortType_AudioBuffer, 0},
		},
	}
}

//-----------------------------------------------------------------------------

// We can't reach the target level with the asymptotic rise/fall of exponentials.
// We will change state when we are within level_epsilon of the target level.
const level_epsilon = 0.001

// Return a k value to give the exponential rise/fall in the required time.
func get_k(t float32, rate int) float32 {
	if t <= 0 {
		return 1.0
	}
	return float32(1.0 - math.Exp(math.Log(level_epsilon)/(float64(t)*float64(rate))))
}

//-----------------------------------------------------------------------------

type adsrState int

const (
	stateNull adsrState = iota
	stateIdle
	stateAttack
	stateDecay
	stateSustain
	stateRelease
)

type adsrModule struct {
	state     adsrState // envelope state
	s         float32   // sustain level
	ka        float32   // attack constant
	kd        float32   // decay constant
	kr        float32   // release constant
	d_trigger float32   // attack->decay trigger level
	s_trigger float32   // decay->sustain trigger level
	i_trigger float32   // release->idle trigger level
	val       float32   // output value
}

// NewADSR returns an Attack/Decay/Sustain/Release envelope module.
func NewADSR() core.Module {
	log.Info.Printf("")
	return &adsrModule{}
}

// Stop stops and performs any cleanup of a module.
func (m *adsrModule) Stop() {
	log.Info.Printf("")
}

//-----------------------------------------------------------------------------
// Events

type adsrEvent int

const (
	null          adsrEvent = iota
	gate                    // attack (!= 0) or release (==0)
	attack_time             // set the attack time
	decay_time              // set the decay time
	sustain_level           // set the sustain level
	release_time            // set the release time
)

func (m *adsrModule) event(etype adsrEvent, val float32) {
	switch etype {
	case gate: // attack (!= 0) or release (==0)
		if val != 0 {
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
	case attack_time: // set the attack time
		if val < 0 {
			panic(fmt.Sprintf("bad attack time %f", val))
		}
		m.ka = get_k(val, core.AUDIO_FS)
	case decay_time: // set the decay time
		if val < 0 {
			panic(fmt.Sprintf("bad decay time %f", val))
		}
		m.kd = get_k(val, core.AUDIO_FS)
	case sustain_level: // set the sustain level
		if val < 0 || val > 1 {
			panic(fmt.Sprintf("bad sustain level %f", val))
		}
		m.s = val
		m.d_trigger = 1.0 - level_epsilon
		m.s_trigger = val + (1.0-val)*level_epsilon
		m.i_trigger = val * level_epsilon
	case release_time: // set the release time
		if val < 0 {
			panic(fmt.Sprintf("bad release time %f", val))
		}
		m.kr = get_k(val, core.AUDIO_FS)
	default:
		panic(fmt.Sprintf("unhandled event type %d", etype))
	}
}

// Event processes a module event.
func (m *adsrModule) Event(e *core.Event) {
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
			if m.val < m.d_trigger {
				m.val += m.ka * (1.0 - m.val)
			} else {
				// goto decay state
				m.val = 1
				m.state = stateDecay
			}
		case stateDecay:
			// decay until sustain level
			if m.val > m.s_trigger {
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
			if m.val > m.i_trigger {
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
