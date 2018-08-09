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

var SimpleInfo = core.PatchInfo{
	Name:   "simple",
	Create: NewSimple,
}

type Simple struct {
	adsr *env.ADSR
	sine *osc.Sine
	pan  *audio.Pan
}

func NewSimple() core.Patch {
	s := &Simple{
		adsr: env.NewADSR(),
		sine: osc.NewSine(),
		pan:  audio.NewPan(),
	}
	return s
}

func (p *Simple) Active() bool {
	return p.adsr.Active()
}

func (p *Simple) Process(in_l, in_r, out_l, out_r *core.SBuf) {
	var env, out core.SBuf
	// generate the envelope
	p.adsr.Process(&env)
	// generate the sine wave
	p.sine.Process(&out)
	// apply the envelope
	core.Mul_SS(&out, &env)
	// pan to left/right channels
	p.pan.Process(&out, out_l, out_r)
}

//-----------------------------------------------------------------------------
