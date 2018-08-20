//-----------------------------------------------------------------------------
/*

Polyphonic Patch

Manage multiple instances of a given subpatch.

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
)

//-----------------------------------------------------------------------------

const MAX_VOICES = 16

//-----------------------------------------------------------------------------

type polyPatch struct {
	subpatch func() core.Patch
	voice    [MAX_VOICES]core.Patch
}

func NewPolyPatch(subpatch func() core.Patch) core.Patch {
	log.Info.Printf("")
	return &polyPatch{
		subpatch: subpatch,
	}
}

//-----------------------------------------------------------------------------

// Run the patch.
func (p *polyPatch) Process(in, out []*core.Buf) {
	for _, v := range p.voice {
		if v != nil && v.Active() {
			v.Process(nil, []*core.Buf{out[0]})
		}
	}
}

// Process a patch event.
func (p *polyPatch) Event(e *core.Event) {
}

// Return true if the patch has non-zero output.
func (p *polyPatch) Active() bool {
	return true
}

//-----------------------------------------------------------------------------
