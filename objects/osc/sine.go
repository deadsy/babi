//-----------------------------------------------------------------------------
/*

Sine Oscillator

*/
//-----------------------------------------------------------------------------

package osc

import "github.com/deadsy/babi/core"

//-----------------------------------------------------------------------------

// frequency to x scaling (xrange/fs)
const FREQ_SCALE = (1 << 32) / core.AUDIO_FS

type Sine struct {
	freq  float32 // base frequency
	x     uint32  // current x-value
	xstep uint32  // current x-step
}

func NewSine() *Sine {
	return &Sine{}
}

func (s *Sine) SetFrequency(freq float32) {
	s.freq = freq
	s.xstep = uint32(s.freq * FREQ_SCALE)
}

func (s *Sine) Process(out *core.Buf) {
	for i := 0; i < len(out); i++ {
		out[i] = core.CosLookup(s.x)
		s.x += s.xstep
	}
}

//-----------------------------------------------------------------------------
