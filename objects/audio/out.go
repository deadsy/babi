//-----------------------------------------------------------------------------
/*

Audio Output Objects

*/
//-----------------------------------------------------------------------------

package audio

import "github.com/deadsy/babi/babi"

//-----------------------------------------------------------------------------
// stereo output

type OutStereo struct {
}

func NewOutStereo() *OutStereo {
	return &OutStereo{}
}

func (o *OutStereo) Process(l, r *babi.SBuf) {
}

//-----------------------------------------------------------------------------
// left/right panning

type Pan struct {
	vol_l, vol_r float32
}

func NewPan() *Pan {
	return &Pan{}
}

func (p *Pan) Process(in, out_l, out_r *babi.SBuf) {
	babi.Copy_SK(out_l, in, p.vol_l)
	babi.Copy_SK(out_r, in, p.vol_r)
}

//-----------------------------------------------------------------------------
