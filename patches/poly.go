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
func (p *polyPatch) Process() {
	for _, v := range p.voice {
		if v != nil && v.Active() {
			v.Process()
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

// Output to the parent patch.
func (p *polyPatch) Out(out ...*core.Buf) {
}

//-----------------------------------------------------------------------------
