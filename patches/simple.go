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
	"github.com/deadsy/babi/objects/noise"
	"github.com/deadsy/babi/objects/osc"
)

//-----------------------------------------------------------------------------

var SimpleInfo = core.PatchInfo{
	Name: "simple",
	New:  NewSimple,
}

type Simple struct {
	adsr  *env.ADSR
	sine  *osc.Sine
	noise *noise.Pink1
	pan   *audio.Pan
	out   *audio.OutLR
}

func NewSimple(b *core.Babi) core.Patch {
	s := &Simple{
		//adsr: env.NewADSR(0.1, 1.0, 0.5, 1.0),
		adsr:  env.NewAD(0.1, 1.0),
		sine:  osc.NewSine(),
		noise: noise.NewPink1(),
		pan:   audio.NewPan(),
		out:   audio.NewOutLR(b),
	}

	s.sine.SetFrequency(440.0)
	s.pan.SetPan(0.5)
	s.pan.SetVol(1.0)
	s.adsr.Attack()

	return s
}

func (p *Simple) Active() bool {
	return p.adsr.Active()
}

func (p *Simple) Process() {
	var env, x0, x1, out_l, out_r core.Buf
	// generate the envelope
	p.adsr.Process(&env)
	// generate the sine wave
	p.sine.Process(&x0)
	// generate noise
	p.noise.Process(&x1)
	// sum the waves
	x0.Add(&x1)
	// apply the envelope
	x0.Mul(&env)
	// pan to left/right channels
	p.pan.Process(&x0, &out_l, &out_r)
	// stereo output
	p.out.Process(&out_l, &out_r)
}

//-----------------------------------------------------------------------------
