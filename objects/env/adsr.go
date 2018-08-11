//-----------------------------------------------------------------------------
/*

ADSR Envelope Object

*/
//-----------------------------------------------------------------------------

package env

import (
	"math"

	"github.com/deadsy/babi/core"
)

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

type ADSRState int

const (
	idle ADSRState = iota
	attack
	decay
	sustain
	release
)

var state_txt = map[ADSRState]string{
	idle:    "idle",
	attack:  "attack",
	decay:   "decay",
	sustain: "sustain",
	release: "release",
}

func (x ADSRState) String() string {
	return state_txt[x]
}

//-----------------------------------------------------------------------------

type ADSR struct {
	s         float32   // sustain level
	ka        float32   // attack constant
	kd        float32   // decay constant
	kr        float32   // release constant
	d_trigger float32   // attack->decay trigger level
	s_trigger float32   // decay->sustain trigger level
	i_trigger float32   // release->idle trigger level
	state     ADSRState // envelope state
	val       float32   // output value
}

// NewADSR returns an Attack/Decay/Sutain/Release envelope generator.
func NewADSR(
	a float32, // attack time in seconds
	d float32, // decay time in seconds
	s float32, // sustain level
	r float32, // release time in seconds
) *ADSR {

	if a < 0 {
		panic("bad attack time")
	}
	if d < 0 {
		panic("bad decay time")
	}
	if s < 0 || s > 1.0 {
		panic("bad sustain level")
	}
	if r < 0 {
		panic("bad release time")
	}

	return &ADSR{
		s:         s,
		ka:        get_k(a, core.AUDIO_FS),
		kd:        get_k(d, core.AUDIO_FS),
		kr:        get_k(r, core.AUDIO_FS),
		d_trigger: 1.0 - level_epsilon,
		s_trigger: s + (1.0-s)*level_epsilon,
		i_trigger: s * level_epsilon,
	}
}

// NewAD returns an Attack/Decay envelope generator.
func NewAD(
	a float32, // attack time in seconds
	d float32, // decay time in seconds
) *ADSR {
	return NewADSR(a, d, 0, 0)
}

//-----------------------------------------------------------------------------

// Attack causes the ADSR state machine to enter the attack segment.
func (e *ADSR) Attack() {
	e.state = attack
}

// Release causes the ADSR state machine to enter the release segment.
func (e *ADSR) Release() {
	if e.state != idle {
		if e.kr == 1 {
			// no release - goto idle
			e.val = 0
			e.state = idle
		} else {
			e.state = release
		}
	}
}

// Idle causes the ADSR state machine to enter the idle state.
func (e *ADSR) Idle() {
	e.val = 0
	e.state = idle
}

// Active returns true if the ADSR state is not idle.
func (e *ADSR) Active() bool {
	return e.state != idle
}

//-----------------------------------------------------------------------------

func (e *ADSR) Process(out *core.Buf) {
	for i := 0; i < len(out); i++ {
		switch e.state {
		case idle:
			// idle - do nothing
		case attack:
			// attack until 1.0 level
			if e.val < e.d_trigger {
				e.val += e.ka * (1.0 - e.val)
			} else {
				// goto decay state
				e.val = 1
				e.state = decay
			}
		case decay:
			// decay until sustain level
			if e.val > e.s_trigger {
				e.val += e.kd * (e.s - e.val)
			} else {
				if e.s != 0 {
					// goto sustain state
					e.val = e.s
					e.state = sustain
				} else {
					// no sustain, goto idle state
					e.val = 0
					e.state = idle
				}
			}
		case sustain:
			// sustain - do nothing
		case release:
			// release until idle level
			if e.val > e.i_trigger {
				e.val += e.kr * (0.0 - e.val)
			} else {
				// goto idle state
				e.val = 0
				e.state = idle
			}
		default:
			panic("bad adsr state")
		}
		out[i] = e.val
	}
}

//-----------------------------------------------------------------------------
