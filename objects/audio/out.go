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
}

func NewOutLR() *OutLR {
	return &OutLR{}
}

func (o *OutLR) Process(l, r *core.Buf) {
}

//-----------------------------------------------------------------------------
