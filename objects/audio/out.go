//-----------------------------------------------------------------------------
/*

Audio Output Objects

*/
//-----------------------------------------------------------------------------

package audio

import "github.com/deadsy/babi/core"

//-----------------------------------------------------------------------------
// stereo output

type OutStereo struct {
}

func NewOutStereo() *OutStereo {
	return &OutStereo{}
}

func (o *OutStereo) Process(l, r *core.SBuf) {
}

//-----------------------------------------------------------------------------
// left/right panning

type Pan struct {
	vol_l, vol_r float32
}

func NewPan() *Pan {
	return &Pan{}
}

func (p *Pan) Process(in, out_l, out_r *core.SBuf) {
	core.Copy_SK(out_l, in, p.vol_l)
	core.Copy_SK(out_r, in, p.vol_r)
}

//-----------------------------------------------------------------------------
