//-----------------------------------------------------------------------------
/*

Karplus Strong Patch

*/
//-----------------------------------------------------------------------------

package patches

import (
	"github.com/deadsy/babi/core"
	"github.com/deadsy/babi/log"
	"github.com/deadsy/babi/objects/osc"
)

//-----------------------------------------------------------------------------

type ksPatch struct {
	ks *osc.KarplusStrong
}

func NewKSPatch() core.Patch {
	log.Info.Printf("")
	return &ksPatch{}
}

//-----------------------------------------------------------------------------

// Run the patch.
func (p *ksPatch) Process(in, out []*core.Buf) {
	// generate the ks wave
	p.ks.Process(out[0])
}

// Process a patch event.
func (p *ksPatch) Event(e *core.Event) {
}

// Return true if the patch has non-zero output.
func (p *ksPatch) Active() bool {
	return true
}

func (p *ksPatch) Stop() {
	// do nothing
}

//-----------------------------------------------------------------------------
