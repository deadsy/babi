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
	vol_l, vol_r float32
}

func NewPan() *Pan {
	return &Pan{}
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
