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
	synth *core.Synth // top-level synth object
}

func NewOutLR(s *core.Synth) *OutLR {
	return &OutLR{
		synth: s,
	}
}

func (o *OutLR) Process(l, r *core.Buf) {
	o.synth.OutLR(l, r)
}

//-----------------------------------------------------------------------------
