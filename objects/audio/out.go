//-----------------------------------------------------------------------------
/*

Audio Output Objects

*/
//-----------------------------------------------------------------------------

package audio

import "github.com/deadsy/babi/core"

//-----------------------------------------------------------------------------
// stereo output

type OutLR struct {
	babi *core.Babi
}

func NewOutLR(b *core.Babi) *OutLR {
	return &OutLR{
		babi: b,
	}
}

func (o *OutLR) Process(l, r *core.Buf) {
	o.babi.OutLR(l, r)
}

//-----------------------------------------------------------------------------
