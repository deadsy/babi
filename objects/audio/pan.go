//-----------------------------------------------------------------------------
/*

Left/Right Panning

*/
//-----------------------------------------------------------------------------

package audio

import "github.com/deadsy/babi/core"

//-----------------------------------------------------------------------------
// left/right panning

type Pan struct {
	vol          float32 // overall volume
	pan          float32 // pan value 0 == left, 1 == right
	vol_l, vol_r float32 // left right channel volume
}

func NewPan() *Pan {
	return &Pan{}
}

func (p *Pan) set() {
	p.vol_l = p.vol * core.Cos(p.pan)
	p.vol_r = p.vol * core.Sin(p.pan)
}

func (p *Pan) SetVol(vol float32) {
	// convert to a linear volume
	p.vol = core.Pow2(core.Clamp(vol, 0, 1)) - 1
	p.set()
}

func (p *Pan) SetPan(pan float32) {
	// Use sin/cos so that l*l + r*r = K (constant power)
	p.pan = core.Clamp(pan, 0, 1) * (core.PI / 2)
	p.set()
}

func (p *Pan) Process(in, out_l, out_r *core.Buf) {
	// left
	out_l.Copy(in)
	out_l.MulScalar(p.vol_l)
	// right
	out_r.Copy(in)
	out_r.MulScalar(p.vol_r)
}

//-----------------------------------------------------------------------------
