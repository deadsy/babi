//-----------------------------------------------------------------------------
/*

A simple patch - Just an envelope on a sine wave.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/objects/audio"
	"github.com/deadsy/babi/objects/env"
	"github.com/deadsy/babi/objects/osc"
)

//-----------------------------------------------------------------------------

type Simple struct {
	adsr *env.ADSR
	sine *osc.Sine
	pan  *audio.Pan
	out  *audio.OutStereo
}

func NewSimple() *Simple {
	s := &Simple{
		adsr: env.NewADSR(),
		sine: osc.NewSine(),
		pan:  audio.NewPan(),
		out:  audio.NewOutStereo(),
	}
	return s
}

func (p *Simple) Process() {
	var env, out, out_l, out_r core.SBuf
	// generate the envelope
	p.adsr.Process(&env)
	// generate the sine wave
	p.sine.Process(&out)
	// apply the envelope
	core.Mul_SS(&out, &env)
	// pan to left/right channels
	p.pan.Process(&out, &out_l, &out_r)
	// stereo output
	p.out.Process(&out_l, &out_r)
}

//-----------------------------------------------------------------------------
