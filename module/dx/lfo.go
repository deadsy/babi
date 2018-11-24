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
			{"wave", "waveform (0..5)", core.PortTypeInt, lfoDxWave},
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
	randState  uint32      // random state for s&h
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

var lfoFrequency = [100]float32{
	0.062506, 0.124815, 0.311474, 0.435381, 0.619784,
	0.744396, 0.930495, 1.116390, 1.284220, 1.496880,
	1.567830, 1.738994, 1.910158, 2.081322, 2.252486,
	2.423650, 2.580668, 2.737686, 2.894704, 3.051722,
	3.208740, 3.366820, 3.524900, 3.682980, 3.841060,
	3.999140, 4.159420, 4.319700, 4.479980, 4.640260,
	4.800540, 4.953584, 5.106628, 5.259672, 5.412716,
	5.565760, 5.724918, 5.884076, 6.043234, 6.202392,
	6.361550, 6.520044, 6.678538, 6.837032, 6.995526,
	7.154020, 7.300500, 7.446980, 7.593460, 7.739940,
	7.886420, 8.020588, 8.154756, 8.288924, 8.423092,
	8.557260, 8.712624, 8.867988, 9.023352, 9.178716,
	9.334080, 9.669644, 10.005208, 10.340772, 10.676336,
	11.011900, 11.963680, 12.915460, 13.867240, 14.819020,
	15.770800, 16.640240, 17.509680, 18.379120, 19.248560,
	20.118000, 21.040700, 21.963400, 22.886100, 23.808800,
	24.731500, 25.759740, 26.787980, 27.816220, 28.844460,
	29.872700, 31.228200, 32.583700, 33.939200, 35.294700,
	36.650200, 37.812480, 38.974760, 40.137040, 41.299320,
	42.461600, 43.639800, 44.818000, 45.996200, 47.174400,
}

func (m *lfoDx) lfoDxSetRate(rate int) {
	rate = core.ClampInt(rate, 0, 99)
	m.xstep = uint32(lfoFrequency[rate] * core.FrequencyScale)
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

//-----------------------------------------------------------------------------
// Port Events

func lfoDxRate(cm core.Module, e *core.Event) {
	m := cm.(*lfoDx)
	m.lfoDxSetRate(e.GetEventInt().Val)
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

// Each waveform ranges from -1.0 to 1.0
// Each waveform is 0 at m.x == 0
func (m *lfoDx) sample() float32 {
	// step the phase
	m.x += m.xstep
	// calculate samples as q8.24
	var sample int32
	switch m.wave {
	case LfoTriangle:
		x := m.x + (1 << 30)
		sample = int32(x >> 6)
		sample ^= -int32(x >> 31)
		sample &= (1 << 25) - 1
		sample -= (1 << 24)
	case LfoSawDown:
		sample = -int32(m.x) >> 7
	case LfoSawUp:
		sample = int32(m.x) >> 7
	case LfoSquare:
		sample = int32(m.x & (1 << 31))
		sample = (sample >> 6) | (1 << 24)
	case LfoSine:
		x := m.x - (1 << 30)
		return core.CosLookup(x)
	case LfoSampleAndHold:
		if m.x < m.xstep {
			// 0..253, cycle length = 128, 64 values with bit 7 = 1
			m.randState = ((m.randState * 179) + 17) & 0xff
		}
		sample = int32(m.randState<<24) >> 7
	}
	// convert q8.24 to float
	return float32(sample) / float32(1<<24)
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
