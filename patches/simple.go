//-----------------------------------------------------------------------------
/*

A simple patch - an AD envelope on a sine wave.

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
	Name: "simple",
	New:  NewSimple,
}

type Simple struct {
	adsr *env.ADSR
	sine *osc.Sine
	pan  *audio.Pan
	out  *audio.OutLR
}

func NewSimple(s *core.Synth) core.Patch {
	p := &Simple{
		adsr: env.NewAD(0.1, 1.0),
		sine: osc.NewSine(),
		pan:  audio.NewPan(),
		out:  audio.NewOutLR(s),
	}

	p.sine.SetFrequency(440.0)
	p.pan.SetPan(0.5)
	p.pan.SetVol(1.0)
	p.adsr.Attack()

	return p
}

func (p *Simple) Start() {
}

func (p *Simple) Stop() {
}

func (p *Simple) Active() bool {
	return p.adsr.Active()
}

func (p *Simple) Process() {
	var env, x0, out_l, out_r core.Buf
	// generate the envelope
	p.adsr.Process(&env)
	// generate the sine wave
	p.sine.Process(&x0)
	// apply the envelope
	x0.Mul(&env)
	// pan to left/right channels
	p.pan.Process(&x0, &out_l, &out_r)
	// stereo output
	p.out.Process(&out_l, &out_r)
}

//-----------------------------------------------------------------------------
