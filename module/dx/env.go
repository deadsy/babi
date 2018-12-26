//-----------------------------------------------------------------------------
/*

DX7 Envelope Generator

https://github.com/mmontag/dx7-synth-js/blob/master/src/envelope-dx7.js
http://wiki.music-synthesizer-for-android.googlecode.com/git/img/env.html

Note:
These reference implementations are similar, but not the same.
The env.html version is probably slower and more accurate.
Both are implemented and are switchable at compile time (envAccurate).

*/
//-----------------------------------------------------------------------------

package dx

import (
	"math"

	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

func init() {
	lutInit()
}

var outputLUT [4096]float32

func lutInit() {
	for i := range outputLUT {
		dB := (float64(i) - 3824.0) * 0.0235
		outputLUT[i] = float32(math.Pow(20.0, (dB / 20.0)))
	}
}

var outputLevel = [100]int{
	0, 5, 9, 13, 17, 20, 23, 25, 27, 29, 31, 33, 35, 37, 39,
	41, 42, 43, 45, 46, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61,
	62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
	81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99,
	100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114,
	115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127,
}

//-----------------------------------------------------------------------------

var envDxInfo = core.ModuleInfo{
	Name: "envDx",
	In: []core.PortInfo{
		{"gate", "envelope gate, attack(>0) or release(=0)", core.PortTypeFloat, envDxGate},
	},
	Out: []core.PortInfo{
		{"out", "output", core.PortTypeAudio, nil},
	},
}

// Info returns the module information.
func (m *envDx) Info() *core.ModuleInfo {
	return &m.info
}

//-----------------------------------------------------------------------------

type envDx struct {
	info           core.ModuleInfo // module info
	levels         *[4]int         // levels for this envelope
	rates          *[4]int         // rates for this envelope
	level          float32         // current level
	targetlevel    float32         // target level
	state          int             // current state
	rising         bool            // rising or falling?
	down           bool            // key state
	idx            int             // incremented every sample
	qr             int
	shift          int
	decayIncrement float32 // decay increment
}

// NewEnv returns an DX7 envelope module.
func NewEnv(s *core.Synth, levels, rates *[4]int) core.Module {
	log.Info.Printf("")
	m := &envDx{
		info:   envDxInfo,
		levels: levels,
		rates:  rates,
	}
	return s.Register(m)
}

// Child returns the child modules of this module.
func (m *envDx) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *envDx) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func envDxGate(cm core.Module, e *core.Event) {
	m := cm.(*envDx)
	gate := e.GetEventFloat().Val
	log.Info.Printf("gate %f", gate)
	if gate != 0 {
		// note on
		m.down = true
		m.idx = 0
		m.advance(0)
	} else {
		// note off
		m.down = false
		m.advance(3)
	}
}

//-----------------------------------------------------------------------------

const envAccurate = true

var envmask = [4][8]int{
	{0, 1, 0, 1, 0, 1, 0, 1},
	{0, 1, 0, 1, 0, 1, 1, 1},
	{0, 1, 1, 1, 0, 1, 1, 1},
	{0, 1, 1, 1, 1, 1, 1, 1},
}

func (m *envDx) envEnable() bool {
	i := m.idx
	if m.shift < 0 {
		sm := (1 << uint(-m.shift)) - 1
		if (i & sm) != sm {
			return false
		}
		i >>= uint(-m.shift)
	}
	return envmask[m.qr&3][i&7] != 0
}

func (m *envDx) attackStep(lev float32) float32 {
	if !m.envEnable() {
		return 0
	}
	slope := 17.0 - (lev / 256.0)
	if m.shift > 0 {
		return slope * m.decayIncrement
	}
	return slope
}

func (m *envDx) decayStep() float32 {
	if !m.envEnable() {
		return 0
	}
	if m.shift > 0 {
		return m.decayIncrement
	}
	return 1.0
}

//-----------------------------------------------------------------------------

// advance moves to the next envelope state.
func (m *envDx) advance(newstate int) {
	m.state = newstate
	if m.state < 4 {
		newlevel := m.levels[m.state]
		m.targetlevel = float32(core.Max(0, (outputLevel[newlevel]<<5)-224))
		m.rising = (m.targetlevel - m.level) > 0
		rateScaling := 0
		m.qr = core.Min(63, rateScaling+((m.rates[m.state]*41)>>6))
		m.shift = (m.qr >> 2) - 11
		m.decayIncrement = core.Pow2(float32(m.shift))
	}
}

// sample generates an envelope sample.
func (m *envDx) sample() float32 {
	if m.state < 3 || (m.state < 4 && !m.down) {
		lev := m.level
		if m.rising {
			if envAccurate {
				lev += m.attackStep(lev)
			} else {
				lev += m.decayIncrement * (2.0 + (m.targetlevel-lev)/256.0)
			}
			if lev >= m.targetlevel {
				lev = m.targetlevel
				m.advance(m.state + 1)
			}
		} else {
			if envAccurate {
				lev -= m.decayStep()
			} else {
				lev -= m.decayIncrement
			}
			if lev <= m.targetlevel {
				lev = m.targetlevel
				m.advance(m.state + 1)
			}
		}
		m.level = lev
	}
	m.idx++
	// Convert DX7 level -> dB -> amplitude
	return outputLUT[int(math.Floor(float64(m.level)))]
}

// Process runs the module DSP.
func (m *envDx) Process(buf ...*core.Buf) bool {
	if m.state >= 4 {
		return false
	}
	out := buf[0]
	for i := range out {
		out[i] = m.sample()
	}
	return true
}

//-----------------------------------------------------------------------------
