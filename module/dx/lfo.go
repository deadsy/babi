//-----------------------------------------------------------------------------
/*

DX7 Low Frequency Oscillator

*/
//-----------------------------------------------------------------------------

package dx

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/utils/log"
)

//-----------------------------------------------------------------------------

// Info returns the module information.
func (m *lfoDx) Info() *core.ModuleInfo {
	return &core.ModuleInfo{
		Name: "lfoDx",
		In: []core.PortInfo{
			{"rate", "rate (0..99)", core.PortTypeInt, lfoDxRate},
			{"delay", "delay (0..99)", core.PortTypeInt, lfoDxDelay},
			{"wave", "wave (0..5)", core.PortTypeInt, lfoDxWave},
			{"sync", "key sync (off/on)", core.PortTypeBool, lfoDxSync},
		},
		Out: []core.PortInfo{
			{"out", "output", core.PortTypeAudio, nil},
		},
	}
}

//-----------------------------------------------------------------------------

type lfoDx struct {
	synth      *core.Synth // top-level synth
	unit       uint32
	wave       lfoWaveType // wave type
	sync       bool        // key sync
	x          uint32      // current x-value
	xstep      uint32      // current x-step
	randState  uint8       // random state for s&h
	delayState uint32
	delayInc   uint32
	delayInc2  uint32
}

// NewLFO returns a DX7 low frequency oscillator module.
func NewLFO(s *core.Synth, cfg *lfoConfig) core.Module {
	log.Info.Printf("")

	m := &lfoDx{
		synth: s,
	}

	n := float64(1 << 6)
	k := float64(1<<32) / (15.5 * 11)
	m.unit = uint32(n*k/float64(core.AudioSampleFrequency) + 0.5)

	log.Info.Printf("unit %d", m.unit)

	if cfg != nil {
		m.wave = cfg.wave
		m.lfoDxSetRate(cfg.speed)
		m.lfoDxSetDelay(cfg.delay)
		m.sync = cfg.sync
	}

	return m
}

// Child returns the child modules of this module.
func (m *lfoDx) Child() []core.Module {
	return nil
}

// Stop performs any cleanup of a module.
func (m *lfoDx) Stop() {
}

//-----------------------------------------------------------------------------
// Port Events

func (m *lfoDx) lfoDxSetRate(rate int) {
	rate = core.ClampInt(rate, 0, 99)
	var sr uint32
	if rate == 0 {
		sr = 1
	} else {
		sr = (165 * uint32(rate)) >> 6
	}
	if sr < 160 {
		sr *= 11
	} else {
		sr *= (11 + ((sr - 160) >> 4))
	}
	m.xstep = m.unit * sr
	freq := float32(core.AudioSampleFrequency) * float32(m.xstep) / float32(1<<32)
	log.Info.Printf("rate %d xstep %d freq %f", rate, m.xstep, freq)
}

func lfoDxRate(cm core.Module, e *core.Event) {
	m := cm.(*lfoDx)
	m.lfoDxSetRate(e.GetEventInt().Val)
}

func (m *lfoDx) lfoDxSetDelay(delay int) {
	delay = core.ClampInt(delay, 0, 99)
	a := uint32(99 - delay)
	if a == 99 {
		m.delayInc = ^uint32(0)
		m.delayInc2 = ^uint32(0)
	} else {
		a = (16 + (a & 15)) << (1 + (a >> 4))
		m.delayInc = m.unit * a
		a &= 0xff80
		a = uint32(core.Max(0x80, int(a)))
		m.delayInc2 = m.unit * a
	}
}

func lfoDxDelay(cm core.Module, e *core.Event) {
	m := cm.(*lfoDx)
	m.lfoDxSetDelay(e.GetEventInt().Val)
}

func lfoDxWave(cm core.Module, e *core.Event) {
	m := cm.(*lfoDx)
	m.wave = lfoWaveType(core.ClampInt(e.GetEventInt().Val, 0, 5))
}

func lfoDxSync(cm core.Module, e *core.Event) {
	m := cm.(*lfoDx)
	m.sync = e.GetEventBool().Val
}

//-----------------------------------------------------------------------------

// Each waveform ranges from 0.0 to 1.0
// Each waveform is 0.5 at m.x == 0

func (m *lfoDx) sample() float32 {

	m.x += m.xstep

	x := int32(1 << 23)

	switch m.wave {
	case LfoTriangle:
		x = int32(m.x >> 7)
		x ^= -(int32(m.x) >> 31)
		x &= (1 << 24) - 1
	case LfoSawDown:
		x = int32((^m.x ^ (1 << 31)) >> 8)
	case LfoSawUp:
		x = int32((m.x ^ (1 << 31)) >> 8)
	case LfoSquare:
		x = int32(((^m.x) >> 7) & (1 << 24))
	case LfoSine:
		//x = (1 << 23) + (Sin::lookup(m.x >> 8) >> 1)
	case LfoSampleAndHold:
		if m.x < m.xstep {
			m.randState = m.randState*179 + 17
		}
		x = int32(m.randState ^ (1 << 7))
		x = (x + 1) << 16
	}

	return float32(x) * (1.0 / float32(1<<24))
}

// Process runs the module DSP.
func (m *lfoDx) Process(buf ...*core.Buf) {
	out := buf[0]
	for i := 0; i < len(out); i++ {
		out[i] = m.sample()
	}
}

// Active returns true if the module has non-zero output.
func (m *lfoDx) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
